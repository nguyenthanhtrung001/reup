package mogo

import (
	"context"
	"fmt"
	"reup/internal/models"
	"reup/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r implRepository) GetExistingVideoIds(ctx context.Context, videoIds []string) (map[string]struct{}, error) {
	col := r.getDouyinVideoCollection()

	existingVideoIds := make(map[string]struct{})

	// Lọc các video theo video_id có trong videoIds
	filter := bson.M{"video_id": bson.M{"$in": videoIds}}
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error fetching existing video ids: %v", err)
	}
	defer cursor.Close(ctx)

	// Duyệt qua các video đã tìm được và thêm vào map
	for cursor.Next(ctx) {
		var video models.BiliVideo
		if err := cursor.Decode(&video); err != nil {
			return nil, fmt.Errorf("error decoding video data: %v", err)
		}
		existingVideoIds[video.VideoID] = struct{}{}
	}

	return existingVideoIds, nil
}

// Implement GetBiliSpaceByDouyinLink
func (r implRepository) GetBiliSpaceByDouyinLink(ctx context.Context, douyinLink string) (*models.BiliSpace, error) {
	col := r.getBiliSpaceCollection()
	var biliSpace models.BiliSpace
	filter := bson.M{"douyin_link": douyinLink}

	err := col.FindOne(ctx, filter).Decode(&biliSpace)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Nếu không tìm thấy, trả về nil
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching BiliSpace: %v", err)
	}

	return &biliSpace, nil
}
