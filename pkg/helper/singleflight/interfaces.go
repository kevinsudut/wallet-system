package singleflight

import "context"

type SingleFlightItf interface {
	DoSingleFlight(ctx context.Context, key string, fn func() (interface{}, error)) (interface{}, error, bool)
}
