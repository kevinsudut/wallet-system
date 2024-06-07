package entity

import "fmt"

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
