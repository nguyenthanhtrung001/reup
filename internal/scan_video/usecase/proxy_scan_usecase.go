package usecase

import (
	"context"
	"reup/internal/models"
)

func (uc implUseCase) DoneAllProxyScan(ctx context.Context) error {
	return uc.repo.DoneAllProxyScan(ctx)
}

func (uc implUseCase) DoneProxyScan(ctx context.Context, proxyIP string) error {
	return uc.repo.DoneProxyScan(ctx, proxyIP)
}

func (uc implUseCase) GetAllProxyScan(ctx context.Context) ([]models.ProxyScan, error) {
	return uc.repo.GetAllProxyScan(ctx)
}

func (uc implUseCase) GetProxyScanRandom(ctx context.Context) (*models.ProxyScan, error) {
	return uc.repo.GetProxyScanRandom(ctx)
}

func (uc implUseCase) InsertProxyScan(ctx context.Context, proxyIP string) error {
	return uc.repo.InsertProxyScan(ctx, proxyIP)
}
