package interceptor

import (
	"errors"
	"testing"
)

func Test_isAuthFree(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "/backend.BackendService/Hello",
			args: args{method: "/backend.BackendService/HelloWorld"},
			want: true,
		},
		{
			name: "/backend.BackendService/SignUp",
			args: args{method: "/backend.BackendService/SignUp"},
			want: true,
		},
		{
			name: "/backend.BackendService/Login",
			args: args{method: "/backend.BackendService/Login"},
			want: true,
		},
		{
			name: "/backend.BackendService/GetArticle",
			args: args{method: "/backend.BackendService/GetArticle"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAuthFree(tt.args.method); got != tt.want {
				t.Errorf("isAuthFree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_recoveryFunc(t *testing.T) {
	type args struct {
		p interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "リカバリー関数",
			args:    args{p: errors.New("some_error")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RecoveryFunc(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("recoveryFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
