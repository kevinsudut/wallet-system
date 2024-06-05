package usecasebalance

type ReadBalanceByUserIdRequest struct {
	UserId string
}

type ReadBalanceByUserIdResponse struct {
	Code    int     `json:"-"`
	Balance float64 `json:"balance"`
}

type TopupBalanceRequest struct {
	UserId string
	Amount float64 `json:"amount"`
}

type TopupBalanceResponse struct {
	Code int `json:"-"`
}

type TransferBalanceRequest struct {
	UserId     string
	ToUsername string  `json:"to_username"`
	Amount     float64 `json:"amount"`
}

type TransferBalanceResponse struct {
	Code int `json:"-"`
}
