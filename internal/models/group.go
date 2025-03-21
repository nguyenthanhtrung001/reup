package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Role      string             `bson:"role"`
	UpdatedAt time.Time          `bson:"updated_at"`
	CreatedAt time.Time          `bson:"created_at"`
}
