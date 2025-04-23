package mogo

import (
	"context"
	"fmt"

	"github.com/nguyenthanhtrung001/reup/internal/models"
	"github.com/nguyenthanhtrung001/reup/internal/scan_video/repository"
	"github.com/nguyenthanhtrung001/reup/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	appSettingCollection = "app_settings"
)

func (repo implRepository) getSettingCollection() mongo.Collection {
	return repo.database.Collection(appSettingCollection)
}

func (repo implRepository) GetFirstRecord(ctx context.Context) (*models.AppSetting, error) {
	col := repo.getSettingCollection()
	filter := bson.D{{}}

	var result models.AppSetting
	err := col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrDocumentNotFound
		}
		return nil, fmt.Errorf("error retrieving document: %v", err)
	}

	return &result, nil
}
