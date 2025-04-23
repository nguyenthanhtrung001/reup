package mongo

import (
	"context"

	"github.com/nguyenthanhtrung001/reup/internal/models"
	"github.com/nguyenthanhtrung001/reup/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	groupCollection = "groups"
)

func (repo implRepository) getGroupCollection() mongo.Collection {
	return repo.database.Collection(groupCollection)
}

func (repo implRepository) GetGroup(ctx context.Context, sc models.Scope, id string) (models.Group, error) {
	col := repo.getGroupCollection()

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		repo.l.Warnf(ctx, "group.repository.mongo.GetGroup.primitive.ObjectIDFromHex: %v", err)
		return models.Group{}, err
	}

	filter := bson.M{"_id": objID}

	var group models.Group
	err = col.FindOne(ctx, filter).Decode(&group)
	if err != nil {
		repo.l.Warnf(ctx, "group.repository.mongo.GetGroup.FindOne: %v", err)
		return models.Group{}, err
	}

	return group, nil
}
