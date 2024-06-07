package domainauth

import (
	"context"

	"github.com/kevinsudut/wallet-system/app/entity"
)

type DomainItf interface {
	InsertUser(ctx context.Context, user entity.User) (err error)
	GetUserById(ctx context.Context, id string) (resp entity.User, err error)
	GetUserByUsername(ctx context.Context, username string) (resp entity.User, err error)
}
