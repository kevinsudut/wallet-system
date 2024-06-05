package usecasebalance

import (
	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
	domainbalance "github.com/kevinsudut/wallet-system/app/domain/balance"
)

type usecase struct {
	balance domainbalance.DomainItf
	auth    domainauth.DomainItf
}

func Init(balance domainbalance.DomainItf, auth domainauth.DomainItf) UsecaseItf {
	return &usecase{
		balance: balance,
		auth:    auth,
	}
}
