package usecaseauth

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
)

func (u usecase) RegisterUser(ctx context.Context, req RegisterUserRequest) (resp RegisterUserResponse, err error) {
	user, err := u.auth.GetUserByUsername(ctx, req.Username)
	if err != nil && err != sql.ErrNoRows {
		return RegisterUserResponse{
			Code: http.StatusBadGateway,
		}, err
	}

	if user.Username != "" {
		return RegisterUserResponse{
			Code: http.StatusConflict,
		}, fmt.Errorf("username already exists")
	}

	err = u.auth.InsertUser(ctx, domainauth.User{
		Username: req.Username,
	})
	if err != nil {
		return RegisterUserResponse{
			Code: http.StatusBadGateway,
		}, err
	}

	token, err := u.token.Create(time.Hour, req.Username)
	if err != nil {
		return RegisterUserResponse{
			Code: http.StatusBadGateway,
		}, err
	}

	return RegisterUserResponse{
		Token: token,
	}, nil
}
