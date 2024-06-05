package handlerauth

import (
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	usecaseauth "github.com/kevinsudut/wallet-system/app/usecase/auth"
	"github.com/kevinsudut/wallet-system/pkg/helper/response"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
)

func (h handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorln("RegisterUser.ReadAll", err)
		response.WriteErrorResponse(w, http.StatusBadRequest)
		return
	}

	var req usecaseauth.RegisterUserRequest

	err = jsoniter.Unmarshal(body, &req)
	if err != nil {
		log.Errorln("RegisterUser.Unmarshal", err)
		response.WriteErrorResponse(w, http.StatusBadRequest)
		return
	}

	resp, err := h.usecase.RegisterUser(r.Context(), req)
	if err != nil {
		log.Errorln("RegisterUser.RegisterUser", err)
		response.WriteErrorResponse(w, resp.Code)
		return
	}

	response.WriteJsonResponse(w, resp.Code, resp)
}
