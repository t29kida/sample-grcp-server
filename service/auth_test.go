package service

import (
	"testing"
)

func Test_auth_CreateAccessToken(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "uuidが生成されること",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAuth()
			got, err := a.CreateAccessToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("CreateAccessToken() got = %v", got)
			}
		})
	}
}
