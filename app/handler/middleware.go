package handler

import (
	"net/http"

	"github.com/kevinsudut/wallet-system/pkg/helper/context"
	"github.com/kevinsudut/wallet-system/pkg/helper/response"
)

var noNeedAuth = map[string]bool{
	"/create_user": true,
}

func (h handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if _, ok := noNeedAuth[r.URL.Path]; !ok {
			data, err := h.token.Validate(r.Header.Get("Authorization"))
			if err != nil {
				response.WriteErrorResponse(w, http.StatusUnauthorized)
				return
			}

			ctx = context.SetAuth(ctx, data.(string))
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
