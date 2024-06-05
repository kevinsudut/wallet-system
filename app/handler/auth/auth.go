package handlerauth

import (
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	usecaseauth "github.com/kevinsudut/wallet-system/app/usecase/auth"
	"github.com/kevinsudut/wallet-system/pkg/helper/response"
)

func (h handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest)
		return
	}

	var req usecaseauth.RegisterUserRequest

	err = jsoniter.Unmarshal(body, &req)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest)
		return
	}

	resp, err := h.usecase.RegisterUser(r.Context(), req)
	if err != nil {
		response.WriteErrorResponse(w, resp.Code)
		return
	}

	response.WriteJsonResponse(w, http.StatusCreated, resp)
}
