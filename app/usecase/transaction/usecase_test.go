package usecasetransaction

import (
	"reflect"
	"testing"

	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
	domainbalance "github.com/kevinsudut/wallet-system/app/domain/balance"
	"github.com/kevinsudut/wallet-system/pkg/helper/singleflight"
)

func TestInit(t *testing.T) {
	type args struct {
		auth    domainauth.DomainItf
		balance domainbalance.DomainItf
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
				balance:      nil,
				auth:         nil,
				singleflight: singleflight.Init(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(tt.args.auth, tt.args.balance); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}
