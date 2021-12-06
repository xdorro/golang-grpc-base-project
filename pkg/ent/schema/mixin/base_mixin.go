// Package mixin entgo
package mixin

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/AmreeshTyagi/goldflake"
	"github.com/spf13/cast"
	"github.com/spf13/viper"

	"github.com/xdorro/golang-grpc-base-project/pkg/ent/hook"
)

// BaseMixin struct.
type BaseMixin struct {
	mixin.Schema
}

// Fields of the BaseMixin.
func (BaseMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").
			Unique(),
	}
}

// Hooks of the BaseMixin.
func (BaseMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		IDHook(),
	}
}

// IDHook create a UUID
func IDHook() ent.Hook {
	gf := goldflake.NewGoldflake(goldflake.Settings{
		StartTime: time.Date(2021, 12, 5, 0, 0, 0, 0, time.UTC),
		MachineID: func() (uint16, error) {
			return cast.ToUint16(viper.Get("MACHINE_ID")), nil
		},
	})

	type IDSetter interface {
		SetID(uint64)
	}

	hk := func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			is, ok := m.(IDSetter)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation %T", m)
			}

			if gf == nil {
				return nil, fmt.Errorf("goldflake not created")
			}

			id, err := gf.NextID()
			if err != nil {
				return nil, fmt.Errorf("id not generated")
			}
			is.SetID(id)
			return next.Mutate(ctx, m)
		})
	}

	return hook.On(hk, ent.OpCreate)
}
