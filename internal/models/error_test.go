package models

import (
	"reflect"
	"testing"
)

func TestError(t *testing.T) {
	type args struct {
		msg  string
		code int
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "normal flow",
			args: args{
				msg:  "Not Found",
				code: 404,
			},
			want: map[string]string{"message": "Not Found", "code": "404"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Error(tt.args.msg, tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorPOSTResponse(t *testing.T) {
	type args struct {
		inputValue string
	}
	tests := []struct {
		name string
		args args
		want *Data
	}{
		{
			name: "normal flow",
			args: args{inputValue: "=+-75"},
			want: &Data{
				Value:  "=+-75",
				Result: "ERROR",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrorPOSTResponse(tt.args.inputValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorPOSTResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
