package handlerbalance

import (
	"fmt"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	usecasebalance "github.com/kevinsudut/wallet-system/app/usecase/balance"
	"github.com/kevinsudut/wallet-system/pkg/helper/context"
	"github.com/kevinsudut/wallet-system/pkg/helper/response"
)

func (h handler) ReadBalance(w http.ResponseWriter, r *http.Request) {
	resp, err := h.usecase.ReadBalanceByUsername(r.Context(), usecasebalance.ReadBalanceByUsernameRequest{
		Username: context.GetAuth(r.Context()),
	})
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest)
		return
	}

	response.WriteJsonResponse(w, http.StatusOK, resp)
}

func (h handler) TopupBalance(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("A", err)
		response.WriteErrorResponse(w, http.StatusBadRequest)
		return
	}

	var req usecasebalance.TopupBalanceRequest

	err = jsoniter.Unmarshal(body, &req)
	if err != nil {
		fmt.Println("B", err)
		response.WriteErrorResponse(w, http.StatusBadRequest)
		return
	}

	req.Username = context.GetAuth(r.Context())

	resp, err := h.usecase.TopupBalance(r.Context(), req)
	if err != nil {
		fmt.Println("C", err)
		response.WriteErrorResponse(w, resp.Code)
		return
	}

	response.WriteJsonResponse(w, http.StatusNoContent, resp)
}

func (h handler) TransferBalance(w http.ResponseWriter, r *http.Request) {

}
