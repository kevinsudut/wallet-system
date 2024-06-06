package usecaseauth

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
	"github.com/kevinsudut/wallet-system/pkg/lib/token"
	gomock "go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	log.Init()
	os.Exit(m.Run())
}

func Test_usecase_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDomainAuth := domainauth.NewMockDomainItf(ctrl)
	mockToken := token.NewMockTokenItf(ctrl)

	type fields struct {
		auth  domainauth.DomainItf
		token token.TokenItf
	}
	type args struct {
		ctx context.Context
		req RegisterUserRequest
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp RegisterUserResponse
		wantErr  bool
		mock     func()
	}{
		{
			name: "success",
			fields: fields{
				auth:  mockDomainAuth,
				token: mockToken,
			},
			args: args{
				ctx: context.Background(),
				req: RegisterUserRequest{
					Username: "username",
				},
			},
			wantResp: RegisterUserResponse{
				Code:  http.StatusCreated,
				Token: "token",
			},
			wantErr: false,
			mock: func() {
				gomock.InOrder(
					mockDomainAuth.EXPECT().GetUserByUsername(gomock.Any(), "username").Return(domainauth.User{}, sql.ErrNoRows),
					mockDomainAuth.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(nil),
					mockToken.EXPECT().Create(time.Hour, gomock.Any()).Return("token", nil),
				)
			},
		},
		{
			name: "error token.Create",
			fields: fields{
				auth:  mockDomainAuth,
				token: mockToken,
			},
			args: args{
				ctx: context.Background(),
				req: RegisterUserRequest{
					Username: "username",
				},
			},
			wantResp: RegisterUserResponse{
				Code: http.StatusBadGateway,
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDomainAuth.EXPECT().GetUserByUsername(gomock.Any(), "username").Return(domainauth.User{}, sql.ErrNoRows),
					mockDomainAuth.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(nil),
					mockToken.EXPECT().Create(time.Hour, gomock.Any()).Return("", fmt.Errorf("foo")),
				)
			},
		},
		{
			name: "error auth.InsertUser",
			fields: fields{
				auth:  mockDomainAuth,
				token: mockToken,
			},
			args: args{
				ctx: context.Background(),
				req: RegisterUserRequest{
					Username: "username",
				},
			},
			wantResp: RegisterUserResponse{
				Code: http.StatusBadGateway,
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDomainAuth.EXPECT().GetUserByUsername(gomock.Any(), "username").Return(domainauth.User{}, sql.ErrNoRows),
					mockDomainAuth.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(fmt.Errorf("foo")),
				)
			},
		},
		{
			name: "error auth.GetUserByUsername",
			fields: fields{
				auth:  mockDomainAuth,
				token: mockToken,
			},
			args: args{
				ctx: context.Background(),
				req: RegisterUserRequest{
					Username: "username",
				},
			},
			wantResp: RegisterUserResponse{
				Code: http.StatusBadGateway,
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDomainAuth.EXPECT().GetUserByUsername(gomock.Any(), "username").Return(domainauth.User{}, fmt.Errorf("foo")),
				)
			},
		},
		{
			name: "error username already exists",
			fields: fields{
				auth:  mockDomainAuth,
				token: mockToken,
			},
			args: args{
				ctx: context.Background(),
				req: RegisterUserRequest{
					Username: "username",
				},
			},
			wantResp: RegisterUserResponse{
				Code: http.StatusConflict,
			},
			wantErr: true,
			mock: func() {
				gomock.InOrder(
					mockDomainAuth.EXPECT().GetUserByUsername(gomock.Any(), "username").Return(domainauth.User{
						Id:       "id",
						Username: "username",
					}, nil),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := usecase{
				auth:  tt.fields.auth,
				token: tt.fields.token,
			}
			tt.mock()
			gotResp, err := u.RegisterUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("usecase.RegisterUser() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
