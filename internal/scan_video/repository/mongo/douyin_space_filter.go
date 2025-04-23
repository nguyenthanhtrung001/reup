package mogo

import (
	"github.com/nguyenthanhtrung001/reup/internal/scan_video/repository"

	"go.mongodb.org/mongo-driver/bson"
)

func (repo implRepository) buildFilter(opts *repository.FindDouyinOldSpacesOptions) bson.M {
	filter := bson.M{
		"type":     "douyin",
		"username": bson.M{"$in": opts.ChannelUsernames},
	}

	// Thêm các điều kiện tùy chỉnh từ Conditions

	return filter
}
