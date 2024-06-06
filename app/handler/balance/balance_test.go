package handlerbalance

import (
	"bytes"
	ctx "context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
	handlertemplate "github.com/kevinsudut/wallet-system/app/handler/template"
	usecasebalance "github.com/kevinsudut/wallet-system/app/usecase/balance"
	"github.com/kevinsudut/wallet-system/pkg/helper/context"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
	"go.uber.org/mock/gomock"
)

func TestMain(t *testing.M) {
	log.Init()
	os.Exit(t.Run())
}

func Test_handler_ReadBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecaseBalance := usecasebalance.NewMockUsecaseItf(ctrl)

	ctx := context.SetAuth(ctx.Background(), domainauth.User{
		Id:       "id",
		Username: "username",
	})

	type fields struct {
		usecase usecasebalance.UsecaseItf
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mock   func()
	}{
		{
			name: "success",
			fields: fields{
				usecase: mockUsecaseBalance,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/balance_read", nil).WithContext(ctx),
			},
			mock: func() {
				gomock.InOrder(
					mockUsecaseBalance.EXPECT().ReadBalanceByUserId(gomock.Any(), usecasebalance.ReadBalanceByUserIdRequest{
						UserId: "id",
					}).Return(usecasebalance.ReadBalanceByUserIdResponse{
						Code:    http.StatusOK,
						Balance: 0,
					}, nil),
				)
			},
		},
		{
			name: "error balance.ReadBalanceByUserId",
			fields: fields{
				usecase: mockUsecaseBalance,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/balance_read", nil).WithContext(ctx),
			},
			mock: func() {
				gomock.InOrder(
					mockUsecaseBalance.EXPECT().ReadBalanceByUserId(gomock.Any(), usecasebalance.ReadBalanceByUserIdRequest{
						UserId: "id",
					}).Return(usecasebalance.ReadBalanceByUserIdResponse{
						Code: http.StatusInternalServerError,
					}, fmt.Errorf("foo")),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler{
				usecase: tt.fields.usecase,
			}
			tt.mock()
			h.ReadBalance(tt.args.w, tt.args.r)
		})
	}
}

func Test_handler_TopupBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecaseBalance := usecasebalance.NewMockUsecaseItf(ctrl)

	ctx := context.SetAuth(ctx.Background(), domainauth.User{
		Id:       "id",
		Username: "username",
	})

	type fields struct {
		usecase usecasebalance.UsecaseItf
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mock   func()
	}{
		{
			name: "success",
			fields: fields{
				usecase: mockUsecaseBalance,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/balance_topup", bytes.NewBufferString(`{"amount":1000}`)).WithContext(ctx),
			},
			mock: func() {
				gomock.InOrder(
					mockUsecaseBalance.EXPECT().TopupBalance(gomock.Any(), usecasebalance.TopupBalanceRequest{
						UserId: "id",
						Amount: 1000,
					}).Return(usecasebalance.TopupBalanceResponse{
						Code: http.StatusNoContent,
					}, nil),
				)
			},
		},
		{
			name: "error balance.TopupBalance",
			fields: fields{
				usecase: mockUsecaseBalance,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/balance_topup", bytes.NewBufferString(`{"amount":1000}`)).WithContext(ctx),
			},
			mock: func() {
				gomock.InOrder(
					mockUsecaseBalance.EXPECT().TopupBalance(gomock.Any(), usecasebalance.TopupBalanceRequest{
						UserId: "id",
						Amount: 1000,
					}).Return(usecasebalance.TopupBalanceResponse{
						Code: http.StatusInternalServerError,
					}, fmt.Errorf("foo")),
				)
			},
		},
		{
			name: "error unmarshal",
			fields: fields{
				usecase: mockUsecaseBalance,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/balance_topup", bytes.NewBufferString(`{"amount":"1000"}`)).WithContext(ctx),
			},
			mock: func() {},
		},
		{
			name: "error read body",
			fields: fields{
				usecase: mockUsecaseBalance,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/balance_topup", handlertemplate.ErrReader{}).WithContext(ctx),
			},
			mock: func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler{
				usecase: tt.fields.usecase,
			}
			tt.mock()
			h.TopupBalance(tt.args.w, tt.args.r)
		})
	}
}

func Test_handler_TransferBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecaseBalance := usecasebalance.NewMockUsecaseItf(ctrl)

	ctx := context.SetAuth(ctx.Background(), domainauth.User{
		Id:       "id",
		Username: "username",
	})

	type fields struct {
		usecase usecasebalance.UsecaseItf
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mock   func()
	}{
		{
			name: "success",
			fields: fields{
				usecase: mockUsecaseBalance,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewBufferString(`{"to_username":"tousername","amount":1000}`)).WithContext(ctx),
			},
			mock: func() {
				gomock.InOrder(
					mockUsecaseBalance.EXPECT().TransferBalance(gomock.Any(), usecasebalance.TransferBalanceRequest{
						UserId:     "id",
						ToUsername: "tousername",
						Amount:     1000,
					}).Return(usecasebalance.TransferBalanceResponse{
						Code: http.StatusNoContent,
					}, nil),
				)
			},
		},
		{
			name: "error balance.TopupBalance",
			fields: fields{
				usecase: mockUsecaseBalance,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewBufferString(`{"to_username":"tousername","amount":1000}`)).WithContext(ctx),
			},
			mock: func() {
				gomock.InOrder(
					mockUsecaseBalance.EXPECT().TransferBalance(gomock.Any(), usecasebalance.TransferBalanceRequest{
						UserId:     "id",
						ToUsername: "tousername",
						Amount:     1000,
					}).Return(usecasebalance.TransferBalanceResponse{
						Code: http.StatusInternalServerError,
					}, fmt.Errorf("foo")),
				)
			},
		},
		{
			name: "error unmarshal",
			fields: fields{
				usecase: mockUsecaseBalance,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewBufferString(`{"amount":"1000"}`)).WithContext(ctx),
			},
			mock: func() {},
		},
		{
			name: "error read body",
			fields: fields{
				usecase: mockUsecaseBalance,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/transfer", handlertemplate.ErrReader{}).WithContext(ctx),
			},
			mock: func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler{
				usecase: tt.fields.usecase,
			}
			tt.mock()
			h.TransferBalance(tt.args.w, tt.args.r)
		})
	}
}
