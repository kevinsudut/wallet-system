package usecaseauth

type RegisterUserRequest struct {
	Username string `json:"username"`
}

type RegisterUserResponse struct {
	Code  int    `json:"-"`
	Token string `json:"token"`
}
