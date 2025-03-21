package usecase

import (
	"context"
	"reup/internal/models"
)

// GetRandomBookTask implements UseCase.
func (uc implUseCase) GetRandomBookTask() (models.Book, error) {
	Book, err := uc.repo.GetRandomBook()
	if err != nil {
		return models.Book{}, err
	}
	if Book.Author != "" {
		uc.l.Info(context.Background(), "chạy gửi pruducer")
		uc.pubRandomBookTask(context.Background(), Book.Title)
	}
	return Book, nil
}
