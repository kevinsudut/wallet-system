package usecaseauth

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kevinsudut/wallet-system/app/entity"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
)

func (u usecase) RegisterUser(ctx context.Context, req RegisterUserRequest) (resp RegisterUserResponse, err error) {
	user, err := u.auth.GetUserByUsername(ctx, req.Username)
	if err != nil && err != sql.ErrNoRows {
		log.Errorln("RegisterUser.GetUserByUsername", err)
		return RegisterUserResponse{
			Code: http.StatusBadGateway,
		}, err
	}

	if user.Id != "" {
		return RegisterUserResponse{
			Code: http.StatusConflict,
		}, fmt.Errorf("username already exists")
	}

	user = entity.User{
		Id:       uuid.NewString(),
		Username: req.Username,
	}

	err = u.auth.InsertUser(ctx, user)
	if err != nil {
		log.Errorln("RegisterUser.InsertUser", err)
		return RegisterUserResponse{
			Code: http.StatusBadGateway,
		}, err
	}

	token, err := u.token.Create(time.Hour, user)
	if err != nil {
		log.Errorln("RegisterUser.Create", err)
		return RegisterUserResponse{
			Code: http.StatusBadGateway,
		}, err
	}

	return RegisterUserResponse{
		Code:  http.StatusCreated,
		Token: token,
	}, nil
}
