package entity

type Balance struct {
	UserId string  `db:"user_id"`
	Amount float64 `db:"amount"`
}
