package handlertemplate

import "github.com/gorilla/mux"

type HandlerItf interface {
	RegisterHandlers(router *mux.Router) *mux.Router
}
