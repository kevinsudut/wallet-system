package usecasetransaction

type ListOverallTopTransactingUsersByValueRequest struct {
	UserId string
}

type ListOverallTopTransactingUsersByValue struct {
	Username        string  `json:"username"`
	TransactedValue float64 `json:"transacted_value"`
}

type ListOverallTopTransactingUsersByValueResponse struct {
	Code int `json:"-"`
	Data []ListOverallTopTransactingUsersByValue
}

type TopTransactionsForUserRequest struct {
	UserId string
}

type TopTransactionsForUser struct {
	Username string  `json:"username"`
	Amount   float64 `json:"amount"`
}

type TopTransactionsForUserResponse struct {
	Code int `json:"-"`
	Data []TopTransactionsForUser
}
