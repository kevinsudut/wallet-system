package handler

import (
	"github.com/gorilla/mux"
)

func (h handler) RegisterHandlers(router *mux.Router) *mux.Router {
	router.Use(h.authMiddleware)

	for _, h := range h.handlers {
		router = h.RegisterHandlers(router)
	}

	return router
}
