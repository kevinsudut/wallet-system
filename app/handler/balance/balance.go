package handlerbalance

import (
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	usecasebalance "github.com/kevinsudut/wallet-system/app/usecase/balance"
	"github.com/kevinsudut/wallet-system/pkg/helper/context"
	"github.com/kevinsudut/wallet-system/pkg/helper/response"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
)

func (h handler) ReadBalance(w http.ResponseWriter, r *http.Request) {
	resp, err := h.usecase.ReadBalanceByUserId(r.Context(), usecasebalance.ReadBalanceByUserIdRequest{
		UserId: context.GetAuth(r.Context()).Id,
	})
	if err != nil {
		log.Errorln("ReadBalance.ReadBalanceByUserId", err)
		response.WriteErrorResponse(w, resp.Code)
		return
	}

	response.WriteJsonResponse(w, resp.Code, resp)
}

func (h handler) TopupBalance(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorln("TopupBalance.ReadAll", err)
		response.WriteErrorResponse(w, http.StatusBadRequest)
		return
	}

	var req usecasebalance.TopupBalanceRequest

	err = jsoniter.Unmarshal(body, &req)
	if err != nil {
		log.Errorln("TopupBalance.Unmarshal", err)
		response.WriteErrorResponse(w, http.StatusBadRequest)
		return
	}

	req.UserId = context.GetAuth(r.Context()).Id

	resp, err := h.usecase.TopupBalance(r.Context(), req)
	if err != nil {
		log.Errorln("TopupBalance.TopupBalance", err)
		response.WriteErrorResponse(w, resp.Code)
		return
	}

	response.WriteJsonResponse(w, resp.Code, resp)
}

func (h handler) TransferBalance(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorln("TransferBalance.ReadAll", err)
		response.WriteErrorResponse(w, http.StatusBadRequest)
		return
	}

	var req usecasebalance.TransferBalanceRequest

	err = jsoniter.Unmarshal(body, &req)
	if err != nil {
		log.Errorln("TransferBalance.Unmarshal", err)
		response.WriteErrorResponse(w, http.StatusBadRequest)
		return
	}

	req.UserId = context.GetAuth(r.Context()).Id

	resp, err := h.usecase.TransferBalance(r.Context(), req)
	if err != nil {
		log.Errorln("TransferBalance.TransferBalance", err)
		response.WriteErrorResponse(w, resp.Code)
		return
	}

	response.WriteJsonResponse(w, resp.Code, resp)
}
