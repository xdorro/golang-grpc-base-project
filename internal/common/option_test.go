package common

import (
	"context"
	"reflect"
	"testing"
)

func TestNewOptions(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		want *Option
	}{
		{
			name: "success",
			want: &Option{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOption(ctx)
			if got != nil {
				got.Ctx = nil
				got.Log = nil
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
