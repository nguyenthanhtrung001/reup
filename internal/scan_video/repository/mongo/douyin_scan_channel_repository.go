package mogo

import (
	"context"
	"fmt"
	"reup/internal/models"
	"reup/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	scanChannelCollection = "douyin_scan_channel"
)

func (repo implRepository) getScanChannelCollection() mongo.Collection {
	return repo.database.Collection(scanChannelCollection)
}
func (repo implRepository) FindDouyinScanChannelBySecUID(ctx context.Context, secUID string) (*models.DouyinScanChannel, error) {

	col := repo.getScanChannelCollection()

	filter := bson.M{"sec_uid": secUID}

	var result models.DouyinScanChannel
	err := col.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding DouyinScanChannel by sec_uid: %v", err)
	}

	return &result, nil
}

func (repo implRepository) InsertDouyinScanChannels(ctx context.Context, channels []models.DouyinScanChannel) error {

	col := repo.getScanChannelCollection()

	var interfaceChannels []interface{}
	for _, channel := range channels {
		interfaceChannels = append(interfaceChannels, channel)
	}

	_, err := col.InsertMany(ctx, interfaceChannels)
	if err != nil {
		return fmt.Errorf("error inserting douyin scan channels: %v", err)
	}

	return nil
}
