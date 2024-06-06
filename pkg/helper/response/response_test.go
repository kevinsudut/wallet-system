package response

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteErrorResponse(t *testing.T) {
	type args struct {
		w          http.ResponseWriter
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				w:          httptest.NewRecorder(),
				statusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteErrorResponse(tt.args.w, tt.args.statusCode)
		})
	}
}

func TestWriteJsonResponse(t *testing.T) {
	type args struct {
		w          http.ResponseWriter
		statusCode int
		content    interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				w:          httptest.NewRecorder(),
				statusCode: 200,
			},
		},
		{
			name: "error",
			args: args{
				w:          httptest.NewRecorder(),
				statusCode: 200,
				content:    make(chan bool),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteJsonResponse(tt.args.w, tt.args.statusCode, tt.args.content)
		})
	}
}
