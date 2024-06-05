package handlerauth

import (
	handlertemplate "github.com/kevinsudut/wallet-system/app/handler/template"
	usecaseauth "github.com/kevinsudut/wallet-system/app/usecase/auth"
)

type handler struct {
	usecase usecaseauth.UsecaseItf
}

func Init(usecase usecaseauth.UsecaseItf) handlertemplate.HandlerItf {
	return &handler{
		usecase: usecase,
	}
}
