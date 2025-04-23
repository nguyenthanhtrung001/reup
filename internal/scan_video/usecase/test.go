package usecase

import (
	"context"
	"reup/internal/models"
)

func (uc implUseCase) GetQuest(ctx context.Context, computer string) (*models.BiliVideo, map[string]interface{}, error) {
	return uc.repo.GetQuest(ctx, computer)
}
