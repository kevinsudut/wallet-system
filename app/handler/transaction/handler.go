package handlertransaction

import (
	handlertemplate "github.com/kevinsudut/wallet-system/app/handler/template"
	usecasetransaction "github.com/kevinsudut/wallet-system/app/usecase/transaction"
)

type handler struct {
	usecase usecasetransaction.UsecaseItf
}

func Init(usecase usecasetransaction.UsecaseItf) handlertemplate.HandlerItf {
	return &handler{
		usecase: usecase,
	}
}
