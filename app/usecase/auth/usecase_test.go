package usecaseauth

import (
	"reflect"
	"testing"

	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
	"github.com/kevinsudut/wallet-system/pkg/lib/token"
)

func TestInit(t *testing.T) {
	type args struct {
		auth  domainauth.DomainItf
		token token.TokenItf
	}
	tests := []struct {
		name string
		args args
		want UsecaseItf
	}{
		{
			args: args{
				auth:  nil,
				token: nil,
			},
			want: &usecase{
				auth:  nil,
				token: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(tt.args.auth, tt.args.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}
