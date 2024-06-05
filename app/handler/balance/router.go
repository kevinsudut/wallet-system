package handlerbalance

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h handler) RegisterHandlers(router *mux.Router) *mux.Router {
	router.HandleFunc("/balance_read", h.ReadBalance).Methods(http.MethodGet)
	router.HandleFunc("/transfer", h.TransferBalance).Methods(http.MethodPost)
	router.HandleFunc("/balance_topup", h.TopupBalance).Methods(http.MethodPost)

	return router
}
