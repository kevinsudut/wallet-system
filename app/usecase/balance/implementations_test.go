package usecasebalance

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
	domainbalance "github.com/kevinsudut/wallet-system/app/domain/balance"
	"github.com/kevinsudut/wallet-system/app/entity"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
	gomock "go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	log.Init()
	os.Exit(m.Run())
}

func Test_usecase_ReadBalanceByUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDomainBalance := domainbalance.NewMockDomainItf(ctrl)

	type fields struct {
		balance domainbalance.DomainItf
		auth    domainauth.DomainItf
	}
	type args struct {
		ctx context.Context
		req ReadBalanceByUserIdRequest
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp ReadBalanceByUserIdResponse
		wantErr  bool
		mock     func()
	}{
		{
			name: "success",
			fields: fields{
				balance: mockDomainBalance,
			},
			args: args{
				ctx: context.Background(),
				req: ReadBalanceByUserIdRequest{
					UserId: "id",
				},
			},
			wantResp: ReadBalanceByUserIdResponse{
				Code:    http.StatusOK,
				Balance: 100,
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GetBalanceByUserId(gomock.Any(), "id").Return(entity.Balance{
						UserId: "id",
						Amount: 100,
					}, nil),
				)
			},
		},
		{
			name: "error balance.GetBalanceByUserId",
			fields: fields{
				balance: mockDomainBalance,
			},
			args: args{
				ctx: context.Background(),
				req: ReadBalanceByUserIdRequest{
					UserId: "id",
				},
			},
			wantResp: ReadBalanceByUserIdResponse{
				Code: http.StatusBadRequest,
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GetBalanceByUserId(gomock.Any(), "id").Return(entity.Balance{}, fmt.Errorf("foo")),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := usecase{
				balance: tt.fields.balance,
				auth:    tt.fields.auth,
			}
			tt.mock()
			gotResp, err := u.ReadBalanceByUserId(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.ReadBalanceByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("usecase.ReadBalanceByUserId() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_usecase_TopupBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDomainBalance := domainbalance.NewMockDomainItf(ctrl)

	type fields struct {
		balance domainbalance.DomainItf
		auth    domainauth.DomainItf
	}
	type args struct {
		ctx context.Context
		req TopupBalanceRequest
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp TopupBalanceResponse
		wantErr  bool
		mock     func()
	}{
		{
			name: "success",
			fields: fields{
				balance: mockDomainBalance,
			},
			args: args{
				ctx: context.Background(),
				req: TopupBalanceRequest{
					UserId: "id",
					Amount: 100,
				},
			},
			wantResp: TopupBalanceResponse{
				Code: http.StatusNoContent,
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GrantBalanceByUserId(gomock.Any(), entity.Balance{
						UserId: "id",
						Amount: 100,
					}).Return(nil),
				)
			},
		},
		{
			name: "error balance.GrantBalanceByUserId",
			fields: fields{
				balance: mockDomainBalance,
			},
			args: args{
				ctx: context.Background(),
				req: TopupBalanceRequest{
					UserId: "id",
					Amount: 100,
				},
			},
			wantResp: TopupBalanceResponse{
				Code: http.StatusBadRequest,
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GrantBalanceByUserId(gomock.Any(), entity.Balance{
						UserId: "id",
						Amount: 100,
					}).Return(fmt.Errorf("foo")),
				)
			},
		},
		{
			name: "error invalid amount",
			fields: fields{
				balance: mockDomainBalance,
			},
			args: args{
				ctx: context.Background(),
				req: TopupBalanceRequest{
					UserId: "id",
					Amount: -100,
				},
			},
			wantResp: TopupBalanceResponse{
				Code: http.StatusBadRequest,
			},
			wantErr: true,
			mock:    func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := usecase{
				balance: tt.fields.balance,
				auth:    tt.fields.auth,
			}
			tt.mock()
			gotResp, err := u.TopupBalance(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.TopupBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("usecase.TopupBalance() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_usecase_TransferBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDomainBalance := domainbalance.NewMockDomainItf(ctrl)
	mockDomainAuth := domainauth.NewMockDomainItf(ctrl)

	type fields struct {
		balance domainbalance.DomainItf
		auth    domainauth.DomainItf
	}
	type args struct {
		ctx context.Context
		req TransferBalanceRequest
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp TransferBalanceResponse
		wantErr  bool
		mock     func()
	}{
		{
			name: "success",
			fields: fields{
				balance: mockDomainBalance,
				auth:    mockDomainAuth,
			},
			args: args{
				ctx: context.Background(),
				req: TransferBalanceRequest{
					UserId:     "id",
					ToUsername: "tousername",
					Amount:     100,
				},
			},
			wantResp: TransferBalanceResponse{
				Code: http.StatusNoContent,
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GetBalanceByUserId(gomock.Any(), "id").Return(entity.Balance{
						UserId: "id",
						Amount: 100,
					}, nil),
					mockDomainAuth.EXPECT().GetUserByUsername(gomock.Any(), "tousername").Return(entity.User{
						Id:       "id",
						Username: "tousername",
					}, nil),
					mockDomainBalance.EXPECT().DisburmentBalance(gomock.Any(), domainbalance.DisburmentBalanceRequest{
						UserId:   "id",
						ToUserId: "id",
						Amount:   100,
					}).Return(nil),
				)
			},
		},
		{
			name: "error balance.DisburmentBalance",
			fields: fields{
				balance: mockDomainBalance,
				auth:    mockDomainAuth,
			},
			args: args{
				ctx: context.Background(),
				req: TransferBalanceRequest{
					UserId:     "id",
					ToUsername: "tousername",
					Amount:     100,
				},
			},
			wantResp: TransferBalanceResponse{
				Code: http.StatusBadRequest,
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GetBalanceByUserId(gomock.Any(), "id").Return(entity.Balance{
						UserId: "id",
						Amount: 100,
					}, nil),
					mockDomainAuth.EXPECT().GetUserByUsername(gomock.Any(), "tousername").Return(entity.User{
						Id:       "id",
						Username: "tousername",
					}, nil),
					mockDomainBalance.EXPECT().DisburmentBalance(gomock.Any(), domainbalance.DisburmentBalanceRequest{
						UserId:   "id",
						ToUserId: "id",
						Amount:   100,
					}).Return(fmt.Errorf("foo")),
				)
			},
		},
		{
			name: "error auth.GetUserByUsername",
			fields: fields{
				balance: mockDomainBalance,
				auth:    mockDomainAuth,
			},
			args: args{
				ctx: context.Background(),
				req: TransferBalanceRequest{
					UserId:     "id",
					ToUsername: "tousername",
					Amount:     100,
				},
			},
			wantResp: TransferBalanceResponse{
				Code: http.StatusNotFound,
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GetBalanceByUserId(gomock.Any(), "id").Return(entity.Balance{
						UserId: "id",
						Amount: 100,
					}, nil),
					mockDomainAuth.EXPECT().GetUserByUsername(gomock.Any(), "tousername").Return(entity.User{
						Id:       "id",
						Username: "tousername",
					}, sql.ErrNoRows),
				)
			},
		},
		{
			name: "error balance.GetBalanceByUserId",
			fields: fields{
				balance: mockDomainBalance,
				auth:    mockDomainAuth,
			},
			args: args{
				ctx: context.Background(),
				req: TransferBalanceRequest{
					UserId:     "id",
					ToUsername: "tousername",
					Amount:     100,
				},
			},
			wantResp: TransferBalanceResponse{
				Code: http.StatusBadRequest,
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GetBalanceByUserId(gomock.Any(), "id").Return(entity.Balance{
						UserId: "id",
						Amount: 100,
					}, fmt.Errorf("foo")),
				)
			},
		},
		{
			name: "error insufficient balance",
			fields: fields{
				balance: mockDomainBalance,
				auth:    mockDomainAuth,
			},
			args: args{
				ctx: context.Background(),
				req: TransferBalanceRequest{
					UserId:     "id",
					ToUsername: "tousername",
					Amount:     100,
				},
			},
			wantResp: TransferBalanceResponse{
				Code: http.StatusBadRequest,
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GetBalanceByUserId(gomock.Any(), "id").Return(entity.Balance{
						UserId: "id",
						Amount: 99,
					}, nil),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := usecase{
				balance: tt.fields.balance,
				auth:    tt.fields.auth,
			}
			tt.mock()
			gotResp, err := u.TransferBalance(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.TransferBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("usecase.TransferBalance() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
