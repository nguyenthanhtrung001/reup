package usecase

import (
	"context"
	"fmt"

	"github.com/nguyenthanhtrung001/reup/internal/scan_video/repository"
)

func (uc implUseCase) GetChannelGroupsByComputer(ctx context.Context) ([]repository.ChannelGroup, error) {
	// Gọi repository để lấy kết quả nhóm theo Computer và đếm Username
	channelGroups, err := uc.repo.FindChannelsGroupedByComputer(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting channel groups by computer: %v", err)
	}

	return channelGroups, nil
}

func (uc implUseCase) CheckChannelExists(ctx context.Context, channelID int64) (bool, error) {
	exists, err := uc.repo.CheckChannelExists(ctx, channelID)
	if err != nil {
		return false, fmt.Errorf("error checking channel existence: %v", err)
	}

	return exists, nil
}
