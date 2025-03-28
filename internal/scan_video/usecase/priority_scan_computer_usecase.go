package usecase

import (
	"context"
	"fmt"
	"reup/internal/models"
)

func (uc implUseCase) FindAllPriorityScanComputers(ctx context.Context) ([]models.PriorityScanComputer, error) {

	computers, err := uc.repo.FindAllPriorityScanComputers(ctx)
	if err != nil {
		return nil, fmt.Errorf("error finding priority scan computers: %v", err)
	}

	return computers, nil
}
