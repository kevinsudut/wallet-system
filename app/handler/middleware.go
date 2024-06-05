package handler

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
	"github.com/kevinsudut/wallet-system/pkg/helper/context"
	"github.com/kevinsudut/wallet-system/pkg/helper/response"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
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
				log.Errorln("authMiddleware.Validate", err)
				response.WriteErrorResponse(w, http.StatusUnauthorized)
				return
			}

			var user domainauth.User
			err = jsoniter.UnmarshalFromString(data.(string), &user)
			if err != nil {
				log.Errorln("authMiddleware.UnmarshalFromString", err)
				response.WriteErrorResponse(w, http.StatusUnauthorized)
				return
			}

			ctx = context.SetAuth(ctx, user)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
