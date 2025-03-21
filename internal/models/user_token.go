package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserTokens struct {
	ID        primitive.ObjectID `bson:"_id"`
	Token     string             `bson:"token"`
	UserID    primitive.ObjectID `bson:"user_id"`
	ExpiredAt time.Time          `bson:"expired_at"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt time.Time          `bson:"deleted_at"`
}
