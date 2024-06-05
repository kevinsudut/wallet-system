package usecasetransaction

import "context"

type UsecaseItf interface {
	ListOverallTopTransactingUsersByValue(ctx context.Context, req ListOverallTopTransactingUsersByValueRequest) (resp ListOverallTopTransactingUsersByValueResponse, err error)
	TopTransactionsForUser(ctx context.Context, req TopTransactionsForUserRequest) (resp TopTransactionsForUserResponse, err error)
}
