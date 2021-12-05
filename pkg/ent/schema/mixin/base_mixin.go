// Package mixin entgo
package mixin

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/oklog/ulid/v2"

	"github.com/kucow/golang-grpc-base/pkg/ent/hook"
)

// BaseMixin struct.
type BaseMixin struct {
	mixin.Schema
}

// Fields of the BaseMixin.
func (BaseMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(35).
			NotEmpty().
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
	type IDSetter interface {
		SetID(string)
	}

	hk := func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			is, ok := m.(IDSetter)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation %T", m)
			}

			t := time.Now().UTC()
			entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
			uuid := ulid.MustNew(ulid.Timestamp(t), entropy)

			is.SetID(uuid.String())
			return next.Mutate(ctx, m)
		})
	}

	return hook.On(hk, ent.OpCreate)
}
