package usecasetransaction

type RegisterUserRequest struct {
	Username string
}

type RegisterUserResponse struct {
	Token string
}
