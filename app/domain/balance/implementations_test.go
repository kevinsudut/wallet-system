package domainbalance

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/kevinsudut/wallet-system/pkg/helper/singleflight"
	"github.com/kevinsudut/wallet-system/pkg/lib/database"
	lrucache "github.com/kevinsudut/wallet-system/pkg/lib/lru-cache"
)

func Test_domain_GetBalanceByUserId(t *testing.T) {
	type fields struct {
		db           database.DatabaseItf
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := domain{
				db:           tt.fields.db,
				cache:        tt.fields.cache,
				stmts:        tt.fields.stmts,
				singleflight: tt.fields.singleflight,
			}
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

func Test_domain_grantBalanceByUserId(t *testing.T) {
	type fields struct {
		db           database.DatabaseItf
		cache        lrucache.LRUCacheItf
		stmts        databaseStmts
		singleflight singleflight.SingleFlightItf
	}
	type args struct {
		ctx     context.Context
		tx      *sql.Tx
		balance Balance
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := domain{
				db:           tt.fields.db,
				cache:        tt.fields.cache,
				stmts:        tt.fields.stmts,
				singleflight: tt.fields.singleflight,
			}
			if err := d.grantBalanceByUserId(tt.args.ctx, tt.args.tx, tt.args.balance); (err != nil) != tt.wantErr {
				t.Errorf("domain.grantBalanceByUserId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_domain_GrantBalanceByUserId(t *testing.T) {
	type fields struct {
		db           database.DatabaseItf
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := domain{
				db:           tt.fields.db,
				cache:        tt.fields.cache,
				stmts:        tt.fields.stmts,
				singleflight: tt.fields.singleflight,
			}
			if err := d.GrantBalanceByUserId(tt.args.ctx, tt.args.balance); (err != nil) != tt.wantErr {
				t.Errorf("domain.GrantBalanceByUserId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_domain_DisburmentBalance(t *testing.T) {
	type fields struct {
		db           database.DatabaseItf
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := domain{
				db:           tt.fields.db,
				cache:        tt.fields.cache,
				stmts:        tt.fields.stmts,
				singleflight: tt.fields.singleflight,
			}
			if err := d.DisburmentBalance(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("domain.DisburmentBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_domain_GetLatestHistoryByUserId(t *testing.T) {
	type fields struct {
		db           database.DatabaseItf
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := domain{
				db:           tt.fields.db,
				cache:        tt.fields.cache,
				stmts:        tt.fields.stmts,
				singleflight: tt.fields.singleflight,
			}
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
	type fields struct {
		db           database.DatabaseItf
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := domain{
				db:           tt.fields.db,
				cache:        tt.fields.cache,
				stmts:        tt.fields.stmts,
				singleflight: tt.fields.singleflight,
			}
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
