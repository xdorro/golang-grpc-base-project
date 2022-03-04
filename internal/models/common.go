package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// IModel is the interface for all models
type IModel interface {
	CollectionName() string
	GetIndexModels() []mongo.IndexModel
}

// Common is the common model
type Common struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// GetIndexModels returns the index models
func (m *Common) GetIndexModels() []mongo.IndexModel {
	return []mongo.IndexModel{}
}

// BeforeCreate is called before creating a new model
func (m *Common) BeforeCreate() {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	m.UpdatedAt = time.Now()
}
