package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title,omitempty"`
	Author      string             `bson:"author,omitempty"`
	PublishedAt string             `bson:"published_at,omitempty"`
	CreateAt    time.Time          `bson:"create_at,omitempty"`
	UpdateAt    time.Time          `bson:"update_at,omitempty"`
}
