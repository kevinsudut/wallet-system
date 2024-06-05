package handler

import (
	handlerauth "github.com/kevinsudut/wallet-system/app/handler/auth"
	handlerbalance "github.com/kevinsudut/wallet-system/app/handler/balance"
	handlertemplate "github.com/kevinsudut/wallet-system/app/handler/template"
	handlertransaction "github.com/kevinsudut/wallet-system/app/handler/transaction"
	"github.com/kevinsudut/wallet-system/app/usecase"
	"github.com/kevinsudut/wallet-system/pkg/lib/database"
	"github.com/kevinsudut/wallet-system/pkg/lib/token"
)

type handler struct {
	handlers []handlertemplate.HandlerItf
	token    token.TokenItf
}

func Init(token token.TokenItf, db database.DatabaseItf) handlertemplate.HandlerItf {
	usecase := usecase.Init(token, db)

	return &handler{
		token: token,
		handlers: []handlertemplate.HandlerItf{
			handlerauth.Init(usecase.Auth),
			handlerbalance.Init(usecase.Balance),
			handlertransaction.Init(usecase.Transaction),
		},
	}
}
