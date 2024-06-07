package context

import (
	"context"

	"github.com/kevinsudut/wallet-system/app/entity"
)

type ctx string

const (
	contextAuth ctx = "context.auth"
)

func SetAuth(ctx context.Context, user entity.User) context.Context {
	return context.WithValue(ctx, contextAuth, user)
}

func GetAuth(ctx context.Context) entity.User {
	user, _ := ctx.Value(contextAuth).(entity.User)
	return user
}
