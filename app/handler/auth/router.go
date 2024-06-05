package handlerauth

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h handler) RegisterHandlers(router *mux.Router) *mux.Router {
	router.HandleFunc("/create_user", h.RegisterUser).Methods(http.MethodPost)

	return router
}
