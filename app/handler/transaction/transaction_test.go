package handlertransaction

import (
	ctx "context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kevinsudut/wallet-system/app/entity"
	usecasetransaction "github.com/kevinsudut/wallet-system/app/usecase/transaction"
	"github.com/kevinsudut/wallet-system/pkg/helper/context"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
	"go.uber.org/mock/gomock"
)

func TestMain(t *testing.M) {
	log.Init()
	os.Exit(t.Run())
}

func Test_handler_ListOverallTopTransactingUsersByValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecaseTransaction := usecasetransaction.NewMockUsecaseItf(ctrl)

	ctx := context.SetAuth(ctx.Background(), entity.User{
		Id:       "id",
		Username: "username",
	})

	type fields struct {
		usecase usecasetransaction.UsecaseItf
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
				usecase: mockUsecaseTransaction,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/top_users", nil).WithContext(ctx),
			},
			mock: func() {
				gomock.InOrder(
					mockUsecaseTransaction.EXPECT().ListOverallTopTransactingUsersByValue(gomock.Any(), usecasetransaction.ListOverallTopTransactingUsersByValueRequest{
						UserId: "id",
					}).Return(usecasetransaction.ListOverallTopTransactingUsersByValueResponse{
						Code: http.StatusOK,
						Data: []usecasetransaction.ListOverallTopTransactingUsersByValue{},
					}, nil),
				)
			},
		},
		{
			name: "error transaction.ListOverallTopTransactingUsersByValue",
			fields: fields{
				usecase: mockUsecaseTransaction,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/top_users", nil).WithContext(ctx),
			},
			mock: func() {
				gomock.InOrder(
					mockUsecaseTransaction.EXPECT().ListOverallTopTransactingUsersByValue(gomock.Any(), usecasetransaction.ListOverallTopTransactingUsersByValueRequest{
						UserId: "id",
					}).Return(usecasetransaction.ListOverallTopTransactingUsersByValueResponse{
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
			h.ListOverallTopTransactingUsersByValue(tt.args.w, tt.args.r)
		})
	}
}

func Test_handler_TopTransactionsForUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecaseTransaction := usecasetransaction.NewMockUsecaseItf(ctrl)

	ctx := context.SetAuth(ctx.Background(), entity.User{
		Id:       "id",
		Username: "username",
	})

	type fields struct {
		usecase usecasetransaction.UsecaseItf
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
				usecase: mockUsecaseTransaction,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/top_transaction_per_user", nil).WithContext(ctx),
			},
			mock: func() {
				gomock.InOrder(
					mockUsecaseTransaction.EXPECT().TopTransactionsForUser(gomock.Any(), usecasetransaction.TopTransactionsForUserRequest{
						UserId: "id",
					}).Return(usecasetransaction.TopTransactionsForUserResponse{
						Code: http.StatusOK,
						Data: []usecasetransaction.TopTransactionsForUser{},
					}, nil),
				)
			},
		},
		{
			name: "error transaction.ListOverallTopTransactingUsersByValue",
			fields: fields{
				usecase: mockUsecaseTransaction,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/top_transaction_per_user", nil).WithContext(ctx),
			},
			mock: func() {
				gomock.InOrder(
					mockUsecaseTransaction.EXPECT().TopTransactionsForUser(gomock.Any(), usecasetransaction.TopTransactionsForUserRequest{
						UserId: "id",
					}).Return(usecasetransaction.TopTransactionsForUserResponse{
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
			h.TopTransactionsForUser(tt.args.w, tt.args.r)
		})
	}
}
