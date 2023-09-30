package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// BusinessCategory represents the category each business_ falls in
type BusinessCategory struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	IsFeatured  bool               `json:"featured" bson:"isFeatured,omitempty"`
	IconUrl     string             `json:"iconUrl" bson:"iconUrl"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
