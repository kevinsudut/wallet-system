package handlerauth

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	handlertemplate "github.com/kevinsudut/wallet-system/app/handler/template"
	usecaseauth "github.com/kevinsudut/wallet-system/app/usecase/auth"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
	"go.uber.org/mock/gomock"
)

func TestMain(t *testing.M) {
	log.Init()
	os.Exit(t.Run())
}

func Test_handler_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecaseAuth := usecaseauth.NewMockUsecaseItf(ctrl)

	type fields struct {
		usecase usecaseauth.UsecaseItf
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
				usecase: mockUsecaseAuth,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/create_user", bytes.NewBufferString(`{"username":"username"}`)),
			},
			mock: func() {
				gomock.InOrder(
					mockUsecaseAuth.EXPECT().RegisterUser(gomock.Any(), usecaseauth.RegisterUserRequest{
						Username: "username",
					}).Return(usecaseauth.RegisterUserResponse{
						Code:  http.StatusCreated,
						Token: "token",
					}, nil),
				)
			},
		},
		{
			name: "error auth.RegisterUser",
			fields: fields{
				usecase: mockUsecaseAuth,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/create_user", bytes.NewBufferString(`{"username":"username"}`)),
			},
			mock: func() {
				gomock.InOrder(
					mockUsecaseAuth.EXPECT().RegisterUser(gomock.Any(), usecaseauth.RegisterUserRequest{
						Username: "username",
					}).Return(usecaseauth.RegisterUserResponse{
						Code: http.StatusInternalServerError,
					}, fmt.Errorf("foo")),
				)
			},
		},
		{
			name: "error unmarshal",
			fields: fields{
				usecase: mockUsecaseAuth,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/create_user", bytes.NewBufferString(`[]`)),
			},
			mock: func() {},
		},
		{
			name: "error read body",
			fields: fields{
				usecase: mockUsecaseAuth,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/create_user", handlertemplate.ErrReader{}),
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
			h.RegisterUser(tt.args.w, tt.args.r)
		})
	}
}
