package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/nguyenthanhtrung001/reup/internal/models"
	"github.com/nguyenthanhtrung001/reup/internal/scan_video/repository"
)

func (uc implUseCase) FindDouyinOldSpaces(ctx context.Context, input *FindDouyinOldSpacesInput) ([]models.BiliSpace, error) {
	opts := &repository.FindDouyinOldSpacesOptions{
		ChannelUsernames:     input.ChannelUsernames,
		ScanNumbers:          input.ScanNumbers,
		DouyinLastScanIsZero: &input.DouyinLastScanIsZero,
	}

	douyinOldSpaces, err := uc.repo.FindDouyinOldSpaces(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding douyin old spaces: %v", err)
	}
	return douyinOldSpaces, nil
}

func (uc implUseCase) CreateBiliSpace(ctx context.Context, input CreateBiliSpaceInput) (*models.BiliSpace, error) {
	// Ánh xạ từ CreateBiliSpaceInput sang repository.CreateBiliSpaceInput
	repoInput := repository.CreateBiliSpaceInput{
		Mid:                input.Mid,
		SpaceTitle:         input.SpaceTitle,
		SpaceThumb:         input.SpaceThumb,
		LastScan:           input.LastScan,
		Note:               input.Note,
		Username:           input.Username,
		CurrentPage:        input.CurrentPage,
		ScanError:          input.ScanError,
		UpdatedAt:          input.UpdatedAt,
		CreatedAt:          time.Now(),
		CountVideo:         input.CountVideo,
		CountUpload:        input.CountUpload,
		CountGet:           input.CountGet,
		CountNext:          input.CountNext,
		LastUpload:         input.LastUpload,
		Type:               input.Type,
		DouyinLink:         input.DouyinLink,
		IxiguaLastScan:     input.IxiguaLastScan,
		FirstJoinScanCount: input.FirstJoinScanCount,
		DouyinLastScan:     input.DouyinLastScan,
		DouyinWaitScan:     input.DouyinWaitScan,
	}

	biliSpace, err := uc.repo.CreateBiliSpace(ctx, repoInput)
	if err != nil {
		return nil, fmt.Errorf("error creating bili space: %v", err)
	}

	return biliSpace, nil
}

func (uc implUseCase) UpdateBiliSpace(ctx context.Context, mid int64) error {
	// Thiết lập các trường cần cập nhật
	updatedFields := map[string]interface{}{
		"douyin_last_scan": time.Now().Unix(),
		"last_scan":        time.Now().Unix(),
		"current_page":     0,
		"scan_error":       0,
	}

	err := uc.repo.UpdateBiliSpaceF(ctx, mid, updatedFields)
	if err != nil {
		return fmt.Errorf("error updating bili space: %v", err)
	}

	return nil
}
