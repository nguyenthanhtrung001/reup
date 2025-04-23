package mogo

import (
	"context"
	"fmt"
	"reup/internal/models"
	"reup/internal/scan_video/repository"
	"reup/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	douyinVideoCollection = "douyin_video"
)

func (repo implRepository) getDouyinVideoCollection() mongo.Collection {
	return repo.database.Collection(douyinVideoCollection)
}
func (repo implRepository) CheckDouyinVideoExists(ctx context.Context, videoID string) bool {
	col := repo.getDouyinVideoCollection()

	filter := bson.M{"video_id": videoID}

	var result models.BiliVideo
	err := col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Nếu không có tài liệu nào, trả về false
			return false
		}
		return false
	}

	return true
}

func (repo implRepository) CreateDouyinVideo(ctx context.Context, input repository.CreateBiliVideoInput) (*models.BiliVideo, error) {
	col := repo.getDouyinVideoCollection()

	// Chuyển đổi từ input sang model BiliVideo
	biliVideo := models.BiliVideo{
		VideoID:           input.VideoID,
		Mid:               input.Mid,
		VideoThumb:        input.VideoThumb,
		UploadTitle:       input.UploadTitle,
		UploadDescription: input.UploadDescription,
		UploadKeyword:     input.UploadKeyword,
		Duration:          input.Duration,
		DateMake:          input.DateMake,
		CountGet:          input.CountGet,
		Next:              input.Next,
		DownloadFail:      input.DownloadFail,
		Type:              input.Type,
		UpdatedAt:         input.UpdatedAt,
		CreatedAt:         input.CreatedAt,
	}

	// Thực hiện insert video mới vào MongoDB
	_, err := col.InsertOne(ctx, biliVideo)
	if err != nil {
		return nil, fmt.Errorf("error inserting bili video: %v", err)
	}

	return &biliVideo, nil
}

func (repo implRepository) InsertBiliVideos(ctx context.Context, videos []models.BiliVideo) error {
	// Kiểm tra nếu slice videos là nil hoặc rỗng
	if videos == nil || len(videos) == 0 {
		return fmt.Errorf("no videos to insert")
	}

	// Lấy collection douyin videos
	col := repo.getDouyinVideoCollection()

	var interfaceVideos []interface{}
	for _, video := range videos {
		// Kiểm tra nếu video là nil hoặc có trường quan trọng không hợp lệ
		if video.Mid == 0 {
			return fmt.Errorf("invalid video data: missing required fields (Mid or DouyinLink)")
		}

		interfaceVideos = append(interfaceVideos, video)
	}

	// Gọi InsertMany với slice interface{}
	_, err := col.InsertMany(ctx, interfaceVideos)
	if err != nil {
		return fmt.Errorf("error inserting bili videos: %v", err)
	}

	return nil
}
