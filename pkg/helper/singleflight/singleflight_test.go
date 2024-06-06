package singleflight

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name string
		want SingleFlightItf
	}{
		{
			want: &singleFlight{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}
