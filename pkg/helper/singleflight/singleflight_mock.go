package singleflight

import "context"

type MockSingleFlight struct{}

func (sf *MockSingleFlight) DoSingleFlight(ctx context.Context, key string, fn func() (interface{}, error)) (interface{}, error, bool) {
	resp, err := fn()
	return resp, err, false
}
