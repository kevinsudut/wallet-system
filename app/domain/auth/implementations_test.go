package domainauth

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kevinsudut/wallet-system/app/entity"
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
	cache.Set(fmt.Sprintf(cacheKeyGetUserById, "id"), entity.User{
		Id:       "id",
		Username: "username",
	}, time.Minute*5)
	cache.Set(fmt.Sprintf(cacheKeyGetUserByUsername, "username"), entity.User{
		Id:       "id",
		Username: "username",
	}, time.Minute*5)
	cache.Set(fmt.Sprintf(cacheKeyGetUserByUsername, "test"), entity.User{}, time.Minute*5)
	os.Exit(m.Run())
}

func Test_domain_InsertUser(t *testing.T) {
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
		ctx  context.Context
		user entity.User
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
					insertUser: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				user: entity.User{
					Id:       "id",
					Username: "username",
				},
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().ExecContextStmt(gomock.Any(), gomock.Any(), "id", "username").Return(nil),
					mockRedis.EXPECT().SetEx(gomock.Any(), fmt.Sprintf(cacheKeyGetUserById, "id"), gomock.Any(), time.Minute*30).Return("", nil),
					mockRedis.EXPECT().SetEx(gomock.Any(), fmt.Sprintf(cacheKeyGetUserByUsername, "username"), gomock.Any(), time.Minute*30).Return("", nil),
					mockCache.EXPECT().Set(fmt.Sprintf(cacheKeyGetUserById, "id"), entity.User{
						Id:       "id",
						Username: "username",
					}, time.Minute*5),
					mockCache.EXPECT().Set(fmt.Sprintf(cacheKeyGetUserByUsername, "username"), entity.User{
						Id:       "id",
						Username: "username",
					}, time.Minute*5),
				)
			},
		},
		{
			name: "error set redis 1",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					insertUser: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				user: entity.User{
					Id:       "id",
					Username: "username",
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().ExecContextStmt(gomock.Any(), gomock.Any(), "id", "username").Return(nil),
					mockRedis.EXPECT().SetEx(gomock.Any(), fmt.Sprintf(cacheKeyGetUserById, "id"), gomock.Any(), time.Minute*30).Return("", fmt.Errorf("foo")),
				)
			},
		},
		{
			name: "error set redis 2",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					insertUser: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				user: entity.User{
					Id:       "id",
					Username: "username",
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().ExecContextStmt(gomock.Any(), gomock.Any(), "id", "username").Return(nil),
					mockRedis.EXPECT().SetEx(gomock.Any(), fmt.Sprintf(cacheKeyGetUserById, "id"), gomock.Any(), time.Minute*30).Return("", nil),
					mockRedis.EXPECT().SetEx(gomock.Any(), fmt.Sprintf(cacheKeyGetUserByUsername, "username"), gomock.Any(), time.Minute*30).Return("", fmt.Errorf("foo")),
				)
			},
		},
		{
			name: "error ExecContextStmt",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					insertUser: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				user: entity.User{
					Id:       "id",
					Username: "username",
				},
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().ExecContextStmt(gomock.Any(), gomock.Any(), "id", "username").Return(fmt.Errorf("foo")),
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
			if err := d.InsertUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("domain.InsertUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_domain_GetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatabase := database.NewMockDatabaseItf(ctrl)
	mockRedis := redis.NewMockRedisItf(ctrl)
	mockCache := lrucache.NewMockLRUCacheItf(ctrl)

	type fields struct {
		db           database.DatabaseItf
		redis        redis.MockRedisItf
		cache        lrucache.LRUCacheItf
		stmts        databaseStmts
		singleflight singleflight.SingleFlightItf
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp entity.User
		wantErr  bool
		mock     func()
	}{
		{
			name: "success",
			fields: fields{
				db:    mockDatabase,
				redis: *mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					getUserById: &sqlx.Stmt{},
				},
				singleflight: &singleflight.MockSingleFlight{},
			},
			args: args{
				ctx: context.Background(),
				id:  "id",
			},
			wantResp: entity.User{
				Id:       "id",
				Username: "username",
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(cacheKeyGetUserById, "id"), time.Minute*5, gomock.Any()).Return(
						cache.Get(fmt.Sprintf(cacheKeyGetUserById, "id")), nil,
					),
				)
			},
		},
		{
			name: "error fetch",
			fields: fields{
				db:    mockDatabase,
				redis: *mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					getUserById: &sqlx.Stmt{},
				},
				singleflight: &singleflight.MockSingleFlight{},
			},
			args: args{
				ctx: context.Background(),
				id:  "id",
			},
			wantResp: entity.User{},
			wantErr:  true,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(cacheKeyGetUserById, "id"), time.Minute*5, gomock.Any()).Return(nil, fmt.Errorf("foo")),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := domain{
				db:           tt.fields.db,
				redis:        &tt.fields.redis,
				cache:        tt.fields.cache,
				stmts:        tt.fields.stmts,
				singleflight: tt.fields.singleflight,
			}
			tt.mock()
			gotResp, err := d.GetUserById(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("domain.GetUserById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("domain.GetUserById() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_domain_GetUserByUsername(t *testing.T) {
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
		ctx      context.Context
		username string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp entity.User
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
					getUserById: &sqlx.Stmt{},
				},
				singleflight: &singleflight.MockSingleFlight{},
			},
			args: args{
				ctx:      context.Background(),
				username: "username",
			},
			wantResp: entity.User{
				Id:       "id",
				Username: "username",
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(cacheKeyGetUserByUsername, "username"), time.Minute*5, gomock.Any()).Return(
						cache.Get(fmt.Sprintf(cacheKeyGetUserByUsername, "username")), nil,
					),
				)
			},
		},
		{
			name: "error no rows",
			fields: fields{
				db:    mockDatabase,
				redis: mockRedis,
				cache: mockCache,
				stmts: databaseStmts{
					getUserById: &sqlx.Stmt{},
				},
				singleflight: &singleflight.MockSingleFlight{},
			},
			args: args{
				ctx:      context.Background(),
				username: "test",
			},
			wantResp: entity.User{},
			wantErr:  true,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(cacheKeyGetUserByUsername, "test"), time.Minute*5, gomock.Any()).Return(cache.Get(fmt.Sprintf(cacheKeyGetUserByUsername, "test")), nil),
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
					getUserById: &sqlx.Stmt{},
				},
				singleflight: &singleflight.MockSingleFlight{},
			},
			args: args{
				ctx:      context.Background(),
				username: "username",
			},
			wantResp: entity.User{},
			wantErr:  true,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(cacheKeyGetUserByUsername, "username"), time.Minute*5, gomock.Any()).Return(nil, fmt.Errorf("foo")),
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
			gotResp, err := d.GetUserByUsername(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("domain.GetUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("domain.GetUserByUsername() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
