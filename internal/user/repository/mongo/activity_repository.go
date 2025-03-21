package mongo

import (
	"context"
	"time"

	"reup/internal/user/repository"
	"reup/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	userActiviyCollection = "user_activities"
)

func (repo implRepository) getUserActitivyCollection() mongo.Collection {
	return repo.database.Collection(userActiviyCollection)
}

func (repo implRepository) CreateActivity(ctx context.Context, opt repository.CreateActivityOptions) error {
	c := repo.getUserActitivyCollection()

	create := bson.M{
		"user_id":    opt.UserID,
		"type":       opt.Type,
		"updated_at": time.Now(),
		"created_at": time.Now(),
	}
	_, err := c.InsertOne(ctx, create)
	if err != nil {
		repo.l.Warnf(ctx, "user.repository.mongo.CreateActivity.InsertOne: %v", err)
		return err
	}

	return nil
}
