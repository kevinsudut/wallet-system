package usecasebalance

import (
	"reflect"
	"testing"

	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
	domainbalance "github.com/kevinsudut/wallet-system/app/domain/balance"
)

func TestInit(t *testing.T) {
	type args struct {
		balance domainbalance.DomainItf
		auth    domainauth.DomainItf
	}
	tests := []struct {
		name string
		args args
		want UsecaseItf
	}{
		{
			args: args{
				balance: nil,
				auth:    nil,
			},
			want: &usecase{
				balance: nil,
				auth:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(tt.args.balance, tt.args.auth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}
