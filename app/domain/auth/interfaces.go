package domainauth

import "context"

type DomainItf interface {
	InsertUser(ctx context.Context, user User) (err error)
	GetUserByUsername(ctx context.Context, username string) (resp User, err error)
}
