package entity

import "github.com/kevinsudut/wallet-system/app/enum"

type History struct {
	Id           string  `db:"id"`
	UserId       string  `db:"user_id"`
	TargetUserId string  `db:"target_user_id"`
	Amount       float64 `db:"amount"`
	Type         int     `db:"type"`
	Notes        string  `db:"notes"`
}

func (h *History) NormalizeAmount() {
	if h.Amount > 0 && h.Type == int(enum.DEBIT) {
		h.Amount *= -1
	}
}
