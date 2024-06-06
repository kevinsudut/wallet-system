package handlerauth

import (
	"reflect"
	"testing"

	handlertemplate "github.com/kevinsudut/wallet-system/app/handler/template"
	usecaseauth "github.com/kevinsudut/wallet-system/app/usecase/auth"
)

func TestInit(t *testing.T) {
	type args struct {
		usecase usecaseauth.UsecaseItf
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
