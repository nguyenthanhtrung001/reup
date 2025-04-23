package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/nguyenthanhtrung001/reup/internal/models"
	"github.com/nguyenthanhtrung001/reup/internal/scan_video/repository"
)

func (uc implUseCase) CheckDouyinVideoExists(ctx context.Context, videoID string) bool {
	exists := uc.repo.CheckDouyinVideoExists(ctx, videoID)
	return exists
}

func (uc implUseCase) CreateDouyinVideo(ctx context.Context, input CreateBiliVideoInput) (*models.BiliVideo, error) {

	repoInput := repository.CreateBiliVideoInput{
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
		Type:              "douyin",
		CreatedAt:         time.Now(),
	}

	// Gọi repository để tạo video mới
	video, err := uc.repo.CreateDouyinVideo(ctx, repoInput)
	if err != nil {
		return nil, fmt.Errorf("error creating bili video: %v", err)
	}

	// Trả về kết quả từ repository
	return video, nil
}
