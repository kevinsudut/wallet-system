package domainbalance

import "fmt"

type Balance struct {
	UserId string  `db:"user_id"`
	Amount float64 `db:"amount"`
}

type DisburmentBalanceRequest struct {
	UserId   string
	ToUserId string
	Amount   float64
}

type History struct {
	Id           string  `db:"id"`
	UserId       string  `db:"user_id"`
	TargetUserId string  `db:"target_user_id"`
	Amount       float64 `db:"amount"`
	Type         int     `db:"type"`
	Notes        string  `db:"notes"`
}

type HistorySummary struct {
	Id           string  `db:"id"`
	UserId       string  `db:"user_id"`
	TargetUserId string  `db:"target_user_id"`
	Amount       float64 `db:"amount"`
	Type         int     `db:"type"`
}

func (hs HistorySummary) GetId() string {
	return fmt.Sprintf("%s:%s:%d", hs.UserId, hs.TargetUserId, hs.Type)
}

func (h *History) NormalizeAmount() {
	if h.Type == int(DEBIT) {
		h.Amount *= -1
	}
}
