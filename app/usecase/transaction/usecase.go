package usecasetransaction

import (
	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
	domainbalance "github.com/kevinsudut/wallet-system/app/domain/balance"
	"github.com/kevinsudut/wallet-system/pkg/helper/singleflight"
)

type usecase struct {
	auth         domainauth.DomainItf
	balance      domainbalance.DomainItf
	singleflight singleflight.SingleFlightItf
}

func Init(auth domainauth.DomainItf, balance domainbalance.DomainItf) UsecaseItf {
	return &usecase{
		auth:         auth,
		balance:      balance,
		singleflight: singleflight.Init(),
	}
}
