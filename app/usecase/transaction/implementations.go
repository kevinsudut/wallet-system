package usecasetransaction

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"sync"

	domainbalance "github.com/kevinsudut/wallet-system/app/domain/balance"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
)

func (u usecase) ListOverallTopTransactingUsersByValue(ctx context.Context, req ListOverallTopTransactingUsersByValueRequest) (resp ListOverallTopTransactingUsersByValueResponse, err error) {
	result, err, _ := u.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyListOverallTopTransactingUsersByValue, req.UserId), func() (interface{}, error) {
		var resp ListOverallTopTransactingUsersByValueResponse

		historySummaries, err := u.balance.GetHistorySummaryByUserIdAndType(ctx, req.UserId, int(domainbalance.DEBIT))
		if err != nil {
			log.Errorln("ListOverallTopTransactingUsersByValue.GetHistorySummaryByUserIdAndType", err)
			return ListOverallTopTransactingUsersByValueResponse{
				Code: http.StatusUnauthorized,
			}, err
		}

		var wg sync.WaitGroup
		goroutineSem := make(chan struct{}, 5)
		errors := make(chan error, len(historySummaries))

		resp.Data = make([]ListOverallTopTransactingUsersByValue, len(historySummaries))
		for idx, historySummary := range historySummaries {
			wg.Add(1)

			go func(idx int, historySummary domainbalance.HistorySummary) {
				goroutineSem <- struct{}{}
				defer func() {
					<-goroutineSem
					wg.Done()
				}()

				user, err := u.auth.GetUserById(ctx, historySummary.TargetUserId)
				if err != nil {
					errors <- err
					log.Errorln("ListOverallTopTransactingUsersByValue.GetUserById", err)
					return
				}

				resp.Data[idx] = ListOverallTopTransactingUsersByValue{
					Username:        user.Username,
					TransactedValue: historySummary.Amount,
				}
			}(idx, historySummary)
		}

		wg.Wait()

		defer func() {
			close(goroutineSem)
			close(errors)
		}()

		if len(errors) > 0 {
			return ListOverallTopTransactingUsersByValueResponse{
				Code: http.StatusUnauthorized,
			}, fmt.Errorf("goroutine error")
		}

		resp.Code = http.StatusOK

		return resp, nil
	})
	if err != nil {
		return result.(ListOverallTopTransactingUsersByValueResponse), err
	}

	return result.(ListOverallTopTransactingUsersByValueResponse), nil
}

func (u usecase) TopTransactionsForUser(ctx context.Context, req TopTransactionsForUserRequest) (resp TopTransactionsForUserResponse, err error) {
	result, err, _ := u.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyListOverallTopTransactingUsersByValue, req.UserId), func() (interface{}, error) {
		var resp TopTransactionsForUserResponse

		histories, err := u.balance.GetLatestHistoryByUserId(ctx, req.UserId)
		if err != nil {
			log.Errorln("TopTransactionsForUser.GetLatestHistoryByUserId", err)
			return TopTransactionsForUserResponse{
				Code: http.StatusUnauthorized,
			}, err
		}

		sort.Slice(histories, func(i, j int) bool {
			return histories[i].Amount > histories[j].Amount
		})

		var wg sync.WaitGroup
		goroutineSem := make(chan struct{}, 5)
		errors := make(chan error, len(histories))

		resp.Data = make([]TopTransactionsForUser, len(histories))
		for idx, history := range histories {
			wg.Add(1)

			go func(idx int, history domainbalance.History) {
				goroutineSem <- struct{}{}
				defer func() {
					<-goroutineSem
					wg.Done()
				}()

				user, err := u.auth.GetUserById(ctx, history.TargetUserId)
				if err != nil {
					errors <- err
					log.Errorln("TopTransactionsForUser.GetUserById", err)
					return
				}

				resp.Data[idx] = TopTransactionsForUser{
					Username: user.Username,
					Amount:   history.Amount,
				}
			}(idx, history)
		}

		wg.Wait()

		defer func() {
			close(goroutineSem)
			close(errors)
		}()

		if len(errors) > 0 {
			return TopTransactionsForUserResponse{
				Code: http.StatusUnauthorized,
			}, fmt.Errorf("goroutine error")
		}

		resp.Code = http.StatusOK

		return resp, nil
	})
	if err != nil {
		return result.(TopTransactionsForUserResponse), err
	}

	return result.(TopTransactionsForUserResponse), nil
}
