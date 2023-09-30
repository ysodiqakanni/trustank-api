package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Business struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	CategoryID  primitive.ObjectID `bson:"category_id"`
	Description string             `bson:"description"`
	Website     string             `bson:"website"`
	Phone       string             `bson:"phone"`
	Email       string             `bson:"email"`

	OwnerId       primitive.ObjectID
	OwnerName     string
	OwnerJobTitle string
	//Reviews       []Review
}
