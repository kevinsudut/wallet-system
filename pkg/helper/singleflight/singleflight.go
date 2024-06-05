package singleflight

import "golang.org/x/sync/singleflight"

type singleFlight struct {
	sf singleflight.Group
}

func Init() SingleFlightItf {
	return &singleFlight{}
}
