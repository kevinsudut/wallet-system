package singleflight

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/sync/singleflight"
)

func Test_singleFlight_DoSingleFlight(t *testing.T) {
	type fields struct {
		sf singleflight.Group
	}
	type args struct {
		ctx context.Context
		key string
		fn  func() (interface{}, error)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		want1   bool
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				sf: singleflight.Group{},
			},
			args: args{
				ctx: context.Background(),
				key: "key",
				fn: func() (interface{}, error) {
					return nil, nil
				},
			},
			want:    nil,
			want1:   false,
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				sf: singleflight.Group{},
			},
			args: args{
				ctx: context.Background(),
				key: "key",
				fn: func() (interface{}, error) {
					return nil, fmt.Errorf("foo")
				},
			},
			want:    nil,
			want1:   false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &singleFlight{
				sf: tt.fields.sf,
			}
			got, err, gotShared := s.DoSingleFlight(tt.args.ctx, tt.args.key, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("singleFlight.DoSingleFlight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("singleFlight.DoSingleFlight() got = %v, want %v", got, tt.want)
			}
			if gotShared != tt.want1 {
				t.Errorf("singleFlight.DoSingleFlight() got1 = %v, want %v", gotShared, tt.want1)
			}
		})
	}
}
