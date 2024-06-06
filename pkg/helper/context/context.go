package context

import (
	"context"

	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
)

type ctx string

const (
	contextAuth ctx = "context.auth"
)

func SetAuth(ctx context.Context, user domainauth.User) context.Context {
	return context.WithValue(ctx, contextAuth, user)
}

func GetAuth(ctx context.Context) domainauth.User {
	user, _ := ctx.Value(contextAuth).(domainauth.User)
	return user
}
