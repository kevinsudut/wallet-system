package domainauth

import (
	"context"
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
	gomock "go.uber.org/mock/gomock"
)

var (
	cache = lrucache.Init()
)

func TestMain(m *testing.M) {
	log.Init()
	cache.Set(fmt.Sprintf(memcacheKeyGetUserById, "id"), User{
		Id:       "id",
		Username: "username",
	}, time.Minute*5)
	cache.Set(fmt.Sprintf(memcacheKeyGetUserByUsername, "username"), User{
		Id:       "id",
		Username: "username",
	}, time.Minute*5)
	cache.Set(fmt.Sprintf(memcacheKeyGetUserByUsername, "test"), User{}, time.Minute*5)
	os.Exit(m.Run())
}

func Test_domain_InsertUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatabase := database.NewMockDatabaseItf(ctrl)
	mockCache := lrucache.NewMockLRUCacheItf(ctrl)

	type fields struct {
		db           database.DatabaseItf
		cache        lrucache.LRUCacheItf
		stmts        databaseStmts
		singleflight singleflight.SingleFlightItf
	}
	type args struct {
		ctx  context.Context
		user User
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
				cache: mockCache,
				stmts: databaseStmts{
					insertUser: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				user: User{
					Id:       "id",
					Username: "username",
				},
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockDatabase.EXPECT().ExecContextStmt(gomock.Any(), gomock.Any(), "id", "username").Return(nil),
					mockCache.EXPECT().Set(fmt.Sprintf(memcacheKeyGetUserById, "id"), User{
						Id:       "id",
						Username: "username",
					}, time.Minute*5),
					mockCache.EXPECT().Set(fmt.Sprintf(memcacheKeyGetUserByUsername, "username"), User{
						Id:       "id",
						Username: "username",
					}, time.Minute*5),
				)
			},
		},
		{
			name: "error ExecContextStmt",
			fields: fields{
				db:    mockDatabase,
				cache: mockCache,
				stmts: databaseStmts{
					insertUser: &sqlx.Stmt{},
				},
			},
			args: args{
				ctx: context.Background(),
				user: User{
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
	mockCache := lrucache.NewMockLRUCacheItf(ctrl)

	type fields struct {
		db           database.DatabaseItf
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
		wantResp User
		wantErr  bool
		mock     func()
	}{
		{
			name: "success",
			fields: fields{
				db:    mockDatabase,
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
			wantResp: User{
				Id:       "id",
				Username: "username",
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(memcacheKeyGetUserById, "id"), time.Minute*5, gomock.Any()).Return(
						cache.Get(fmt.Sprintf(memcacheKeyGetUserById, "id")), nil,
					),
				)
			},
		},
		{
			name: "error fetch",
			fields: fields{
				db:    mockDatabase,
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
			wantResp: User{},
			wantErr:  true,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(memcacheKeyGetUserById, "id"), time.Minute*5, gomock.Any()).Return(nil, fmt.Errorf("foo")),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := domain{
				db:           tt.fields.db,
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
	mockCache := lrucache.NewMockLRUCacheItf(ctrl)

	type fields struct {
		db           database.DatabaseItf
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
		wantResp User
		wantErr  bool
		mock     func()
	}{
		{
			name: "success",
			fields: fields{
				db:    mockDatabase,
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
			wantResp: User{
				Id:       "id",
				Username: "username",
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(memcacheKeyGetUserByUsername, "username"), time.Minute*5, gomock.Any()).Return(
						cache.Get(fmt.Sprintf(memcacheKeyGetUserByUsername, "username")), nil,
					),
				)
			},
		},
		{
			name: "error no rows",
			fields: fields{
				db:    mockDatabase,
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
			wantResp: User{},
			wantErr:  true,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(memcacheKeyGetUserByUsername, "test"), time.Minute*5, gomock.Any()).Return(cache.Get(fmt.Sprintf(memcacheKeyGetUserByUsername, "test")), nil),
				)
			},
		},
		{
			name: "error fetch",
			fields: fields{
				db:    mockDatabase,
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
			wantResp: User{},
			wantErr:  true,
			mock: func() {
				gomock.InOrder(
					mockCache.EXPECT().Fetch(fmt.Sprintf(memcacheKeyGetUserByUsername, "username"), time.Minute*5, gomock.Any()).Return(nil, fmt.Errorf("foo")),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := domain{
				db:           tt.fields.db,
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
