package handlerbalance

import (
	"reflect"
	"testing"

	handlertemplate "github.com/kevinsudut/wallet-system/app/handler/template"
	usecasebalance "github.com/kevinsudut/wallet-system/app/usecase/balance"
)

func TestInit(t *testing.T) {
	type args struct {
		usecase usecasebalance.UsecaseItf
	}
	tests := []struct {
		name string
		args args
		want handlertemplate.HandlerItf
	}{
		{
			args: args{
				usecase: nil,
			},
			want: &handler{
				usecase: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(tt.args.usecase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}
