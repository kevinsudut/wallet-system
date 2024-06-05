package handlerbalance

import (
	handlertemplate "github.com/kevinsudut/wallet-system/app/handler/template"
	usecasebalance "github.com/kevinsudut/wallet-system/app/usecase/balance"
)

type handler struct {
	usecase usecasebalance.UsecaseItf
}

func Init(usecase usecasebalance.UsecaseItf) handlertemplate.HandlerItf {
	return &handler{
		usecase: usecase,
	}
}
