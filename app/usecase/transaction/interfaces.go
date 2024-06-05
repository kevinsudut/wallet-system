package usecasetransaction

import "context"

type UsecaseItf interface {
	RegisterUser(ctx context.Context, req RegisterUserRequest) (resp RegisterUserResponse, err error)
}
