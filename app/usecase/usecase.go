package usecase

import (
	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
	domainbalance "github.com/kevinsudut/wallet-system/app/domain/balance"
	usecaseauth "github.com/kevinsudut/wallet-system/app/usecase/auth"
	usecasebalance "github.com/kevinsudut/wallet-system/app/usecase/balance"
	usecasetransaction "github.com/kevinsudut/wallet-system/app/usecase/transaction"
	"github.com/kevinsudut/wallet-system/pkg/lib/database"
	"github.com/kevinsudut/wallet-system/pkg/lib/token"
)

type usecase struct {
	Auth        usecaseauth.UsecaseItf
	Balance     usecasebalance.UsecaseItf
	Transaction usecasetransaction.UsecaseItf
}

func Init(token token.TokenItf, db database.DatabaseItf) usecase {
	domainAuth := domainauth.Init(db)
	domainBalance := domainbalance.Init(db)

	return usecase{
		Auth:        usecaseauth.Init(domainAuth, token),
		Balance:     usecasebalance.Init(domainBalance, domainAuth),
		Transaction: usecasetransaction.Init(domainAuth, domainBalance),
	}
}
