package repo

import (
	"context"
	"reflect"
	"testing"

	"github.com/xdorro/golang-grpc-base-project/internal/common"
)

func TestNewRepo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	opts := common.NewOption(ctx)

	type args struct {
		opts *common.Option
	}
	tests := []struct {
		name string
		args args
		want *Repo
	}{
		{
			name: "success",
			args: args{
				opts: opts,
			},
			want: &Repo{
				Ctx:    opts.Ctx,
				Log:    opts.Log,
				Client: opts.Client,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRepo(tt.args.opts.Ctx, tt.args.opts.Log, tt.args.opts.Client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}
