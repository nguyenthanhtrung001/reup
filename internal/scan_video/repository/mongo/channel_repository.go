package mogo

import (
	"context"
	"fmt"

	"github.com/nguyenthanhtrung001/reup/internal/models"
	"github.com/nguyenthanhtrung001/reup/internal/scan_video/repository"
	"github.com/nguyenthanhtrung001/reup/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	mon "go.mongodb.org/mongo-driver/mongo"
)

const (
	channnelCollection = "channel"
)

func (repo implRepository) getChannelCollection() mongo.Collection {
	return repo.database.Collection(channnelCollection)
}

func (repo implRepository) FindChannelsGroupedByComputer(ctx context.Context) ([]repository.ChannelGroup, error) {
	col := repo.getChannelCollection()

	// Sử dụng aggregation để nhóm theo trường "computer" và đếm số lượng "username"
	pipeline := mon.Pipeline{
		// Bước group: nhóm theo trường "computer" và đếm số lượng "username"
		{{Key: "$group", Value: bson.M{
			"_id":        "$computer",       // Nhóm theo trường "computer"
			"total_user": bson.M{"$sum": 1}, // Đếm số lượng username trong mỗi nhóm
		}}},

		// Bước match: lọc chỉ lấy những nhóm có computer chứa _DY hoặc _XIT
		{{Key: "$match", Value: bson.M{
			"_id": bson.M{
				"$regex": "(_DY|_XIT)$", // Lọc các computer có chứa _DY hoặc _XIT ở cuối
			},
		}}},

		// Bước sort: sắp xếp theo số lượng user giảm dần
		{{Key: "$sort", Value: bson.M{
			"total_user": -1, // Sắp xếp theo số lượng user giảm dần
		}}},
	}

	// Thực hiện aggregation
	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("could not aggregate channel data: %v", err)
	}
	defer cursor.Close(ctx)

	// Duyệt qua các kết quả trả về
	var channelGroups []repository.ChannelGroup
	for cursor.Next(ctx) {
		var group repository.ChannelGroup
		if err := cursor.Decode(&group); err != nil {
			return nil, fmt.Errorf("cursor decode error: %v", err)
		}
		channelGroups = append(channelGroups, group)
	}

	return channelGroups, nil
}

func (repo implRepository) CheckChannelExists(ctx context.Context, channelID int64) (bool, error) {
	col := repo.getChannelCollection()

	filter := bson.M{"id": channelID}

	var result models.Channel
	err := col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Nếu không có tài liệu nào, trả về false
			return false, nil
		}
		return false, fmt.Errorf("error checking channel existence: %v", err)
	}

	return true, nil
}

func (repo implRepository) FindUsernamesByComputer(ctx context.Context, computerName string) ([]string, error) {
	col := repo.getChannelCollection()

	filter := bson.M{"computer": computerName}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("could not find channels for computer %s: %v", computerName, err)
	}
	defer cursor.Close(ctx)

	var usernames []string

	for cursor.Next(ctx) {
		var channel models.Channel
		if err := cursor.Decode(&channel); err != nil {
			return nil, fmt.Errorf("cursor decode error: %v", err)
		}
		usernames = append(usernames, channel.Username) // Thêm username vào danh sách
	}

	return usernames, nil
}
