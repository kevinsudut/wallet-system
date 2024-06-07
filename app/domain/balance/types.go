package domainbalance

type DisburmentBalanceRequest struct {
	UserId   string
	ToUserId string
	Amount   float64
}
