package domainbalance

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kevinsudut/wallet-system/pkg/helper/singleflight"
	"github.com/kevinsudut/wallet-system/pkg/lib/database"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
	lrucache "github.com/kevinsudut/wallet-system/pkg/lib/lru-cache"
	"github.com/kevinsudut/wallet-system/pkg/lib/redis"
	gomock "go.uber.org/mock/gomock"
)

var (
	cache = lrucache.Init()
)

func TestMain(m *testing.M) {
	log.Init()
	cache.Set(fmt.Sprintf(cacheKeyGetBalanceByUserId, "test"), Balance{}, time.Minute*5)
	cache.Set(fmt.Sprintf(cacheKeyGetBalanceByUserId, "id"), Balance{
		UserId: "id",
		Amount: 10,
	}, time.Minute*5)
	cache.Set(fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, "id"), []History{
		{
			UserId:       "id",
			TargetUserId: "id",
			Amount:       10,
			Type:         1,
		},
	}, time.Minute*5)
	cache.Set(fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "id", 1), []HistorySummary{
		{
			UserId:       "id",
			TargetUserId: "id",
			Amount:       10,
			Type:         1,
		},
	}, time.Minute*5)
	os.Exit(m.Run())
}

func Test_domain_GetBalanceByUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatabase := database.NewMockDatabaseItf(ctrl)
	mockRedis := redis.NewMockRedisItf(ctrl)
	mockCache := lrucache.NewMockLRUCacheItf(ctrl)

	type fields struct {
		db           database.DatabaseItf
		redis        redis.RedisItf
		cache        lrucache.LRUCacheItf
		stmts        databaseStmts
		singleflight singleflight.SingleFlightItf
	}
	type args struct {
		ctx    context.Context
		userId string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp Balance
		wantErr  bool
		mock     func()
	}{
		{
			name: "success",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					getBalanceByUserId: &sqlx.Stmt{},
				},
				singleflight: &singleflight.MockSingleFlight{},
			},
			args: args{
				ctx:    context.Background(),
				userId: "id",
			},
			wantResp: Balance{
				UserId: "id",
				Amount: 10,
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(cacheKeyGetBalanceByUserId, "id"), time.Minute*5, gomock.Any()).Return(
						cache.Get(fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")), nil,
					),
				)
			},
		},
		{
			name: "error fetch",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					getBalanceByUserId: &sqlx.Stmt{},
				},
				singleflight: &singleflight.MockSingleFlight{},
			},
			args: args{
				ctx:    context.Background(),
				userId: "id",
			},
			wantResp: Balance{},
			wantErr:  true,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(cacheKeyGetBalanceByUserId, "id"), time.Minute*5, gomock.Any()).Return(nil, fmt.Errorf("foo")),
				)
			},
		},
		{
			name: "error not found",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					getBalanceByUserId: &sqlx.Stmt{},
				},
				singleflight: &singleflight.MockSingleFlight{},
			},
			args: args{
				ctx:    context.Background(),
				userId: "test",
			},
			wantResp: Balance{},
			wantErr:  true,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(cacheKeyGetBalanceByUserId, "test"), time.Minute*5, gomock.Any()).Return(
						cache.Get(fmt.Sprintf(cacheKeyGetBalanceByUserId, "test")), nil,
					),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := domain{
				db:           tt.fields.db,
				redis:        tt.fields.redis,
				cache:        tt.fields.cache,
				stmts:        tt.fields.stmts,
				singleflight: tt.fields.singleflight,
			}
			tt.mock()
			gotResp, err := d.GetBalanceByUserId(tt.args.ctx, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("domain.GetBalanceByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("domain.GetBalanceByUserId() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_domain_GrantBalanceByUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatabase := database.NewMockDatabaseItf(ctrl)
	mockRedis := redis.NewMockRedisItf(ctrl)
	mockCache := lrucache.NewMockLRUCacheItf(ctrl)

	type fields struct {
		db           database.DatabaseItf
		redis        redis.RedisItf
		cache        lrucache.LRUCacheItf
		stmts        databaseStmts
		singleflight singleflight.SingleFlightItf
	}
	type args struct {
		ctx     context.Context
		balance Balance
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				balance: Balance{
					UserId: "id",
					Amount: 10,
				},
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, nil),
					// grantBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(false),

					// insertHistory
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),

					// updateHistorySummary
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "id", 1)).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "id", 1)).Return(false),

					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, "id")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, "id")).Return(false),
					mockDatabase.EXPECT().Commit(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "error insertHistory.Delete redis",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				balance: Balance{
					UserId: "id",
					Amount: 10,
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, nil),
					// grantBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(false),

					// insertHistory
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),

					// updateHistorySummary
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "id", 1)).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "id", 1)).Return(false),

					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, "id")).Return(int64(0), fmt.Errorf("foo")),
					mockDatabase.EXPECT().Rollback(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "error updateHistorySummary.Delete redis",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				balance: Balance{
					UserId: "id",
					Amount: 10,
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, nil),
					// grantBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(false),

					// insertHistory
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),

					// updateHistorySummary
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "id", 1)).Return(int64(0), fmt.Errorf("foo")),
					mockDatabase.EXPECT().Rollback(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "error updateHistorySummary.ExecContextStmtTx db",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				balance: Balance{
					UserId: "id",
					Amount: 10,
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, nil),
					// grantBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(false),

					// insertHistory
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),

					// updateHistorySummary
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("foo")),
					mockDatabase.EXPECT().Rollback(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "error insertHistory.ExecContextStmtTx db",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				balance: Balance{
					UserId: "id",
					Amount: 10,
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, nil),
					// grantBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(false),

					// insertHistory
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("foo")),
					mockDatabase.EXPECT().Rollback(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "error grantBalanceByUserId.Delete redis",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				balance: Balance{
					UserId: "id",
					Amount: 10,
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, nil),
					// grantBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(int64(0), fmt.Errorf("foo")),
					mockDatabase.EXPECT().Rollback(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "error grantBalanceByUserId.ExecContextStmtTx db",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				balance: Balance{
					UserId: "id",
					Amount: 10,
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, nil),
					// grantBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("foo")),
					mockDatabase.EXPECT().Rollback(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "error Begin",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				balance: Balance{
					UserId: "id",
					Amount: 10,
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, fmt.Errorf("foo")),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := domain{
				db:           tt.fields.db,
				redis:        tt.fields.redis,
				cache:        tt.fields.cache,
				stmts:        tt.fields.stmts,
				singleflight: tt.fields.singleflight,
			}
			tt.mock()
			if err := d.GrantBalanceByUserId(tt.args.ctx, tt.args.balance); (err != nil) != tt.wantErr {
				t.Errorf("domain.GrantBalanceByUserId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_domain_DisburmentBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatabase := database.NewMockDatabaseItf(ctrl)
	mockRedis := redis.NewMockRedisItf(ctrl)
	mockCache := lrucache.NewMockLRUCacheItf(ctrl)

	type fields struct {
		db           database.DatabaseItf
		redis        redis.RedisItf
		cache        lrucache.LRUCacheItf
		stmts        databaseStmts
		singleflight singleflight.SingleFlightItf
	}
	type args struct {
		ctx context.Context
		req DisburmentBalanceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				req: DisburmentBalanceRequest{
					UserId:   "id",
					ToUserId: "toid",
					Amount:   10,
				},
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, nil),
					// grantBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "toid")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, "toid")).Return(false),

					// deductBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(false),

					// insertHistory
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "toid", 1)).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "toid", 1)).Return(false),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, "toid")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, "toid")).Return(false),

					// insertHistory
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "id", 2)).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "id", 2)).Return(false),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, "id")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, "id")).Return(false),
					mockDatabase.EXPECT().Commit(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "error insertHistory.ExecContextStmtTx db",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				req: DisburmentBalanceRequest{
					UserId:   "id",
					ToUserId: "toid",
					Amount:   10,
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, nil),
					// grantBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "toid")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, "toid")).Return(false),

					// deductBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(false),

					// insertHistory
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "toid", 1)).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "toid", 1)).Return(false),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, "toid")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, "toid")).Return(false),

					// insertHistory
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("foo")),
					mockDatabase.EXPECT().Rollback(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "error insertHistory.ExecContextStmtTx db",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				req: DisburmentBalanceRequest{
					UserId:   "id",
					ToUserId: "toid",
					Amount:   10,
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, nil),
					// grantBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "toid")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, "toid")).Return(false),

					// deductBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(false),

					// insertHistory
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("foo")),
					mockDatabase.EXPECT().Rollback(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "error deductBalanceByUserId.Delete redis",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				req: DisburmentBalanceRequest{
					UserId:   "id",
					ToUserId: "toid",
					Amount:   10,
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, nil),
					// grantBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "toid")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, "toid")).Return(false),

					// deductBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "id")).Return(int64(0), fmt.Errorf("foo")),
					mockDatabase.EXPECT().Rollback(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "error deductBalanceByUserId.ExecContextStmtTx db",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				req: DisburmentBalanceRequest{
					UserId:   "id",
					ToUserId: "toid",
					Amount:   10,
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, nil),
					// grantBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					mockRedis.EXPECT().Delete(gomock.Any(), fmt.Sprintf(cacheKeyGetBalanceByUserId, "toid")).Return(int64(0), nil),
					mockCache.EXPECT().Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, "toid")).Return(false),

					// deductBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("foo")),
					mockDatabase.EXPECT().Rollback(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "error grantBalanceByUserId.ExecContextStmtTx db",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				req: DisburmentBalanceRequest{
					UserId:   "id",
					ToUserId: "toid",
					Amount:   10,
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, nil),
					// grantBalanceByUserId
					mockDatabase.EXPECT().ExecContextStmtTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("foo")),
					mockDatabase.EXPECT().Rollback(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "error Begin",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					grantBalanceByUserId:     &sqlx.Stmt{},
					insertHistory:            &sqlx.Stmt{},
					updateHistorySummaryById: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				req: DisburmentBalanceRequest{
					UserId:   "id",
					ToUserId: "toid",
					Amount:   10,
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().Begin().Return(&sql.Tx{}, fmt.Errorf("foo")),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := domain{
				db:           tt.fields.db,
				redis:        tt.fields.redis,
				cache:        tt.fields.cache,
				stmts:        tt.fields.stmts,
				singleflight: tt.fields.singleflight,
			}
			tt.mock()
			if err := d.DisburmentBalance(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("domain.DisburmentBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_domain_GetLatestHistoryByUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatabase := database.NewMockDatabaseItf(ctrl)
	mockRedis := redis.NewMockRedisItf(ctrl)
	mockCache := lrucache.NewMockLRUCacheItf(ctrl)

	type fields struct {
		db           database.DatabaseItf
		redis        redis.RedisItf
		cache        lrucache.LRUCacheItf
		stmts        databaseStmts
		singleflight singleflight.SingleFlightItf
	}
	type args struct {
		ctx    context.Context
		userId string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp []History
		wantErr  bool
		mock     func()
	}{
		{
			name: "success",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					getLatestHistoryByUserId: &sqlx.Stmt{},
				},
				singleflight: &singleflight.MockSingleFlight{},
			},
			args: args{
				ctx:    context.Background(),
				userId: "id",
			},
			wantResp: []History{
				{
					UserId:       "id",
					TargetUserId: "id",
					Amount:       10,
					Type:         1,
				},
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, "id"), time.Minute*5, gomock.Any()).Return(
						cache.Get(fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, "id")), nil,
					),
				)
			},
		},
		{
			name: "error fetch",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					getLatestHistoryByUserId: &sqlx.Stmt{},
				},
				singleflight: &singleflight.MockSingleFlight{},
			},
			args: args{
				ctx:    context.Background(),
				userId: "id",
			},
			wantResp: nil,
			wantErr:  true,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, "id"), time.Minute*5, gomock.Any()).Return(nil, fmt.Errorf("foo")),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := domain{
				db:           tt.fields.db,
				redis:        tt.fields.redis,
				cache:        tt.fields.cache,
				stmts:        tt.fields.stmts,
				singleflight: tt.fields.singleflight,
			}
			tt.mock()
			gotResp, err := d.GetLatestHistoryByUserId(tt.args.ctx, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("domain.GetLatestHistoryByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("domain.GetLatestHistoryByUserId() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_domain_GetHistorySummaryByUserIdAndType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatabase := database.NewMockDatabaseItf(ctrl)
	mockRedis := redis.NewMockRedisItf(ctrl)
	mockCache := lrucache.NewMockLRUCacheItf(ctrl)

	type fields struct {
		db           database.DatabaseItf
		redis        redis.RedisItf
		cache        lrucache.LRUCacheItf
		stmts        databaseStmts
		singleflight singleflight.SingleFlightItf
	}
	type args struct {
		ctx         context.Context
		userId      string
		historyType int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp []HistorySummary
		wantErr  bool
		mock     func()
	}{
		{
			name: "success",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					getHistorySummaryByUserIdAndType: &sqlx.Stmt{},
				},
				singleflight: &singleflight.MockSingleFlight{},
			},
			args: args{
				ctx:         context.Background(),
				userId:      "id",
				historyType: 1,
			},
			wantResp: []HistorySummary{
				{
					UserId:       "id",
					TargetUserId: "id",
					Amount:       10,
					Type:         1,
				},
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "id", 1), time.Minute*5, gomock.Any()).Return(
						cache.Get(fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "id", 1)), nil,
					),
				)
			},
		},
		{
			name: "error fetch",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					getHistorySummaryByUserIdAndType: &sqlx.Stmt{},
				},
				singleflight: &singleflight.MockSingleFlight{},
			},
			args: args{
				ctx:         context.Background(),
				userId:      "id",
				historyType: 1,
			},
			wantResp: nil,
			wantErr:  true,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, "id", 1), time.Minute*5, gomock.Any()).Return(nil, fmt.Errorf("foo")),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := domain{
				db:           tt.fields.db,
				redis:        tt.fields.redis,
				cache:        tt.fields.cache,
				stmts:        tt.fields.stmts,
				singleflight: tt.fields.singleflight,
			}
			tt.mock()
			gotResp, err := d.GetHistorySummaryByUserIdAndType(tt.args.ctx, tt.args.userId, tt.args.historyType)
			if (err != nil) != tt.wantErr {
				t.Errorf("domain.GetHistorySummaryByUserIdAndType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("domain.GetHistorySummaryByUserIdAndType() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
