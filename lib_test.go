package spoon_test

import (
	"testing"

	"github.com/pi9min/spoon"
)

func TestQuote(t *testing.T) {
	type args struct {
		unquoted string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "a-z0-9",
			args: args{
				unquoted: "abc123",
			},
			want: "`abc123`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := spoon.Quote(tt.args.unquoted); got != tt.want {
				t.Errorf("Quote() = %v, expect %v", got, tt.want)
			}
		})
	}
}

func TestSemicolon(t *testing.T) {
	type args struct {
		schema string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "a-z0-9",
			args: args{
				schema: "abc123",
			},
			want: "abc123;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := spoon.Semicolon(tt.args.schema); got != tt.want {
				t.Errorf("Semicolon() = %v, expect %v", got, tt.want)
			}
		})
	}
}
