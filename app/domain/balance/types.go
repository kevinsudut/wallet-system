package domainbalance

type Balance struct {
	Username string  `db:"username"`
	Amount   float64 `db:"amount"`
}
