package usecaseauth

import (
	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
	"github.com/kevinsudut/wallet-system/pkg/lib/token"
)

type usecase struct {
	auth  domainauth.DomainItf
	token token.TokenItf
}

func Init(auth domainauth.DomainItf, token token.TokenItf) UsecaseItf {
	return &usecase{
		auth:  auth,
		token: token,
	}
}
