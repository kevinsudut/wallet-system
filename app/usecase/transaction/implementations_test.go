package usecasetransaction

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
	domainbalance "github.com/kevinsudut/wallet-system/app/domain/balance"
	"github.com/kevinsudut/wallet-system/app/entity"
	"github.com/kevinsudut/wallet-system/pkg/helper/singleflight"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
	gomock "go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	log.Init()
	os.Exit(m.Run())
}

func Test_usecase_ListOverallTopTransactingUsersByValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDomainBalance := domainbalance.NewMockDomainItf(ctrl)
	mockDomainAuth := domainauth.NewMockDomainItf(ctrl)
	mockSingleFlight := &singleflight.MockSingleFlight{}

	type fields struct {
		auth         domainauth.DomainItf
		balance      domainbalance.DomainItf
		singleflight singleflight.SingleFlightItf
	}
	type args struct {
		ctx context.Context
		req ListOverallTopTransactingUsersByValueRequest
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp ListOverallTopTransactingUsersByValueResponse
		wantErr  bool
		mock     func()
	}{
		{
			name: "success",
			fields: fields{
				balance:      mockDomainBalance,
				auth:         mockDomainAuth,
				singleflight: mockSingleFlight,
			},
			args: args{
				ctx: context.Background(),
				req: ListOverallTopTransactingUsersByValueRequest{
					UserId: "id",
				},
			},
			wantResp: ListOverallTopTransactingUsersByValueResponse{
				Code: http.StatusOK,
				Data: []ListOverallTopTransactingUsersByValue{
					{
						Username:        "user1",
						TransactedValue: 100,
					},
					{
						Username:        "user1",
						TransactedValue: 200,
					},
				},
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GetHistorySummaryByUserIdAndType(gomock.Any(), "id", 2).Return([]entity.HistorySummary{
						{
							UserId:       "id",
							TargetUserId: "target1",
							Amount:       100,
							Type:         2,
						},
						{
							UserId:       "id",
							TargetUserId: "target1",
							Amount:       200,
							Type:         2,
						},
					}, nil),
					mockDomainAuth.EXPECT().GetUserById(gomock.Any(), "target1").Return(entity.User{
						Id:       "target1",
						Username: "user1",
					}, nil),
					mockDomainAuth.EXPECT().GetUserById(gomock.Any(), "target1").Return(entity.User{
						Id:       "target1",
						Username: "user1",
					}, nil),
				)
			},
		},
		{
			name: "error auth.GetUserById",
			fields: fields{
				balance:      mockDomainBalance,
				auth:         mockDomainAuth,
				singleflight: mockSingleFlight,
			},
			args: args{
				ctx: context.Background(),
				req: ListOverallTopTransactingUsersByValueRequest{
					UserId: "id",
				},
			},
			wantResp: ListOverallTopTransactingUsersByValueResponse{
				Code: http.StatusUnauthorized,
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GetHistorySummaryByUserIdAndType(gomock.Any(), "id", 2).Return([]entity.HistorySummary{
						{
							UserId:       "id",
							TargetUserId: "target1",
							Amount:       100,
							Type:         2,
						},
						{
							UserId:       "id",
							TargetUserId: "target1",
							Amount:       200,
							Type:         2,
						},
					}, nil),
					mockDomainAuth.EXPECT().GetUserById(gomock.Any(), "target1").Return(entity.User{
						Id:       "target1",
						Username: "user1",
					}, nil),
					mockDomainAuth.EXPECT().GetUserById(gomock.Any(), "target1").Return(entity.User{}, fmt.Errorf("foo")),
				)
			},
		},
		{
			name: "error balance.GetHistorySummaryByUserIdAndType",
			fields: fields{
				balance:      mockDomainBalance,
				auth:         mockDomainAuth,
				singleflight: mockSingleFlight,
			},
			args: args{
				ctx: context.Background(),
				req: ListOverallTopTransactingUsersByValueRequest{
					UserId: "id",
				},
			},
			wantResp: ListOverallTopTransactingUsersByValueResponse{
				Code: http.StatusUnauthorized,
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GetHistorySummaryByUserIdAndType(gomock.Any(), "id", 2).Return([]entity.HistorySummary{}, fmt.Errorf("foo")),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := usecase{
				auth:         tt.fields.auth,
				balance:      tt.fields.balance,
				singleflight: tt.fields.singleflight,
			}
			tt.mock()
			gotResp, err := u.ListOverallTopTransactingUsersByValue(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.ListOverallTopTransactingUsersByValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("usecase.ListOverallTopTransactingUsersByValue() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_usecase_TopTransactionsForUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDomainBalance := domainbalance.NewMockDomainItf(ctrl)
	mockDomainAuth := domainauth.NewMockDomainItf(ctrl)
	mockSingleFlight := &singleflight.MockSingleFlight{}

	type fields struct {
		auth         domainauth.DomainItf
		balance      domainbalance.DomainItf
		singleflight singleflight.SingleFlightItf
	}
	type args struct {
		ctx context.Context
		req TopTransactionsForUserRequest
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp TopTransactionsForUserResponse
		wantErr  bool
		mock     func()
	}{
		{
			name: "success",
			fields: fields{
				auth:         mockDomainAuth,
				balance:      mockDomainBalance,
				singleflight: mockSingleFlight,
			},
			args: args{
				ctx: context.Background(),
				req: TopTransactionsForUserRequest{
					UserId: "id",
				},
			},
			wantResp: TopTransactionsForUserResponse{
				Code: http.StatusOK,
				Data: []TopTransactionsForUser{
					{
						Username: "username",
						Amount:   100,
					},
					{
						Username: "username",
						Amount:   50,
					},
				},
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GetLatestHistoryByUserId(gomock.Any(), "id").Return([]entity.History{
						{
							UserId:       "id",
							TargetUserId: "id",
							Amount:       100,
							Type:         1,
						},
						{
							UserId:       "id",
							TargetUserId: "id",
							Amount:       50,
							Type:         1,
						},
					}, nil),
					mockDomainAuth.EXPECT().GetUserById(gomock.Any(), "id").Return(entity.User{
						Id:       "id",
						Username: "username",
					}, nil),
					mockDomainAuth.EXPECT().GetUserById(gomock.Any(), "id").Return(entity.User{
						Id:       "id",
						Username: "username",
					}, nil),
				)
			},
		},
		{
			name: "error auth.GetUserById",
			fields: fields{
				auth:         mockDomainAuth,
				balance:      mockDomainBalance,
				singleflight: mockSingleFlight,
			},
			args: args{
				ctx: context.Background(),
				req: TopTransactionsForUserRequest{
					UserId: "id",
				},
			},
			wantResp: TopTransactionsForUserResponse{
				Code: http.StatusUnauthorized,
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GetLatestHistoryByUserId(gomock.Any(), "id").Return([]entity.History{
						{
							UserId:       "id",
							TargetUserId: "id",
							Amount:       50,
							Type:         1,
						},
					}, nil),
					mockDomainAuth.EXPECT().GetUserById(gomock.Any(), "id").Return(entity.User{}, fmt.Errorf("foo")),
				)
			},
		},
		{
			name: "error auth.GetLatestHistoryByUserId",
			fields: fields{
				auth:         mockDomainAuth,
				balance:      mockDomainBalance,
				singleflight: mockSingleFlight,
			},
			args: args{
				ctx: context.Background(),
				req: TopTransactionsForUserRequest{
					UserId: "id",
				},
			},
			wantResp: TopTransactionsForUserResponse{
				Code: http.StatusUnauthorized,
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDomainBalance.EXPECT().GetLatestHistoryByUserId(gomock.Any(), "id").Return([]entity.History{}, fmt.Errorf("foo")),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := usecase{
				auth:         tt.fields.auth,
				balance:      tt.fields.balance,
				singleflight: tt.fields.singleflight,
			}
			tt.mock()
			gotResp, err := u.TopTransactionsForUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.TopTransactionsForUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("usecase.TopTransactionsForUser() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
