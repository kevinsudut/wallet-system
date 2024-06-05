package handlertransaction

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h handler) RegisterHandlers(router *mux.Router) *mux.Router {
	router.HandleFunc("/top_users", h.ListOverallTopTransactingUsersByValue).Methods(http.MethodGet)
	router.HandleFunc("/top_transaction_per_user", h.TopTransactionsForUser).Methods(http.MethodGet)

	return router
}
