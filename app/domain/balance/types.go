package domainbalance

import "fmt"

type Balance struct {
	Username string  `db:"username"`
	Amount   float64 `db:"amount"`
}

type DisburmentBalanceRequest struct {
	Username   string
	ToUsername string
	Amount     float64
}

type History struct {
	Id             string  `db:"id"`
	Username       string  `db:"username"`
	TargetUsername string  `db:"target_username"`
	Amount         float64 `db:"amount"`
	Type           int     `db:"type"`
	Notes          string  `db:"notes"`
}

type HistorySummary struct {
	Id             string  `db:"id"`
	Username       string  `db:"username"`
	TargetUsername string  `db:"target_username"`
	Amount         float64 `db:"amount"`
	Type           int     `db:"type"`
}

func (hs HistorySummary) GetId() string {
	return fmt.Sprintf("%s:%s:%d", hs.Username, hs.TargetUsername, hs.Type)
}

func (h *History) NormalizeAmount() {
	if h.Type == int(DEBIT) {
		h.Amount *= -1
	}
}
