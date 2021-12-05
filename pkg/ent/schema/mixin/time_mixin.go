// Package mixin entgo
package mixin

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// TimeMixin struct.
type TimeMixin struct {
	mixin.Schema
}

// Fields of the TimeMixin.
func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("create_time").
			Immutable().
			Default(time.Now),

		field.Time("update_time").
			Default(time.Now).
			UpdateDefault(time.Now).
			Immutable(),

		field.Time("delete_time").
			Optional(),
	}
}

// Indexes of the TimeMixin.
func (TimeMixin) Indexes() []ent.Index {
	return []ent.Index{
		// non-unique index.
		index.Fields("create_time", "update_time", "delete_time"),
	}
}
