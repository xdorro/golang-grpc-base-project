package common

import (
	"strings"
	"testing"
)

func TestGenerateFromPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				password: "test_password",
			},
			want: "$2a$10$",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateFromPassword(tt.args.password)
			if err != nil {
				t.Errorf("GenerateFromPassword() error = %v", err)
				return
			}
			if !strings.HasPrefix(got, tt.want) {
				t.Errorf("GenerateFromPassword() got = %v, want prefix = %v", got, tt.want)
			}
		})
	}
}

func TestCompareHashAndPassword(t *testing.T) {
	type args struct {
		hash     string
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success",
			args: args{
				password: "test_password",
				hash:     "$2a$10$QBGPLhR2F6dY3ZYCVkZgRu8j/011qp3Clj4bWr5sy5g8460kCfKfq",
			},
			want: true,
		},
		{
			name: "failed",
			args: args{
				password: "test_password123",
				hash:     "$2a$10$QBGPLhR2F6dY3ZYCVkZgRu8j/011qp3Clj4bWr5sy5g8460kCfKfq",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareHashAndPassword(tt.args.hash, tt.args.password); got != tt.want {
				t.Errorf("CompareHashAndPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
