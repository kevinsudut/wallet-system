package context

import "context"

type ctx string

const (
	contextAuth ctx = "context.auth"
)

func SetAuth(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, contextAuth, username)
}

func GetAuth(ctx context.Context) string {
	return ctx.Value(contextAuth).(string)
}
