package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User represents a user.
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Name           string             `bson:"name"`
	Email          string             `bson:"email"`
	Role           []string           `bson:"role"`
	HashedPassword []byte             `bson:"hashed_password"`
	Created        time.Time          `bson:"created"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

func (u User) GetRole() []string {
	return u.Role
}

// GetID returns the user ID.
func (u User) GetID() primitive.ObjectID {
	return u.ID
}

// GetName returns the user name.
func (u User) GetName() string {
	return u.Email
}
