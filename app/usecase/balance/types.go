package usecasebalance

type ReadBalanceByUsernameRequest struct {
	Username string
}

type ReadBalanceByUsernameResponse struct {
	Balance float64 `json:"balance"`
}

type TopupBalanceRequest struct {
	Username string
	Amount   float64 `json:"amount"`
}

type TopupBalanceResponse struct {
	Code int
}

type TransferBalanceRequest struct {
	Username   string
	ToUsername string  `json:"to_username"`
	Amount     float64 `json:"amount"`
}

type TransferBalanceResponse struct {
	Code int
}
