package repo

import (
	"context"
	"reflect"
	"testing"
)

func TestNewRepo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	opts := optioncommon.NewOption(ctx)

	type args struct {
		opts *optioncommon.Option
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
				ctx:    opts.Ctx,
				log:    opts.Log,
				client: opts.Client,
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
