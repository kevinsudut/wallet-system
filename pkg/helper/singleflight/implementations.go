package singleflight

import "context"

func (s *singleFlight) DoSingleFlight(ctx context.Context, key string, fn func() (interface{}, error)) (interface{}, error, bool) {
	resp, err, shared := s.sf.Do(key, fn)
	if err != nil {
		s.sf.Forget(key)
		return resp, err, shared
	}
	return resp, nil, shared
}
