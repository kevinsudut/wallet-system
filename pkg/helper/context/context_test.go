package context

import (
	"context"
	"reflect"
	"testing"

	domainauth "github.com/kevinsudut/wallet-system/app/domain/auth"
)

func TestSetAuth(t *testing.T) {
	type args struct {
		ctx  context.Context
		user domainauth.User
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			args: args{
				ctx:  context.Background(),
				user: domainauth.User{},
			},
			want: context.WithValue(context.Background(), contextAuth, domainauth.User{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetAuth(tt.args.ctx, tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAuth(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want domainauth.User
	}{
		{
			args: args{
				ctx: context.Background(),
			},
			want: domainauth.User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAuth(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}
