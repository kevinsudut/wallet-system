package handlertransaction

import (
	"net/http"

	usecasetransaction "github.com/kevinsudut/wallet-system/app/usecase/transaction"
	"github.com/kevinsudut/wallet-system/pkg/helper/context"
	"github.com/kevinsudut/wallet-system/pkg/helper/response"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
)

func (h handler) ListOverallTopTransactingUsersByValue(w http.ResponseWriter, r *http.Request) {
	resp, err := h.usecase.ListOverallTopTransactingUsersByValue(r.Context(), usecasetransaction.ListOverallTopTransactingUsersByValueRequest{
		UserId: context.GetAuth(r.Context()).Id,
	})
	if err != nil {
		log.Errorln("ListOverallTopTransactingUsersByValue.ListOverallTopTransactingUsersByValue", err)
		response.WriteErrorResponse(w, resp.Code)
		return
	}

	response.WriteJsonResponse(w, resp.Code, resp.Data)
}

func (h handler) TopTransactionsForUser(w http.ResponseWriter, r *http.Request) {
	resp, err := h.usecase.TopTransactionsForUser(r.Context(), usecasetransaction.TopTransactionsForUserRequest{
		UserId: context.GetAuth(r.Context()).Id,
	})
	if err != nil {
		log.Errorln("TopTransactionsForUser.TopTransactionsForUser", err)
		response.WriteErrorResponse(w, resp.Code)
		return
	}

	response.WriteJsonResponse(w, resp.Code, resp.Data)
}
