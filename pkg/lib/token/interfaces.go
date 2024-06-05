package token

import "time"

type TokenItf interface {
	Create(ttl time.Duration, content interface{}) (string, error)
	Validate(token string) (interface{}, error)
}
