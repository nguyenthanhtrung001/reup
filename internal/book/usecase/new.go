package usecase

import (
	"context"
	prod "reup/internal/book/delivery/rabbitmq/producer"
	"reup/internal/book/repository"
	"reup/internal/models"
	"reup/pkg/encrypter"
	"reup/pkg/log"
)

type UseCase interface {
	EncrypterUC
	BookUseCase
}
type BookUseCase interface {
	CreateBook(ctx context.Context, inputBook BookInput) error
	GetBookBy(ctx context.Context, id string) (BookOutput, error)
	DeleteBook(ctx context.Context, id string) error
	UpdateBook(ctx context.Context, inputUpdate UpdateInput) error
	List(ctx context.Context, sc models.Scope, input ListBookInput) (ListBookOutput, error)
	GetRandomBookTask() (models.Book, error)
}

type EncrypterUC interface {
	Create(ctx context.Context, text string) (string, error)
}
type implUseCase struct {
	l         log.Logger
	repo      repository.Repository
	encrypter encrypter.Encrypter
	prod      prod.Producer
}

func New(l log.Logger, repo repository.Repository, encrypter encrypter.Encrypter, prod prod.Producer) UseCase {
	return implUseCase{
		l:         l,
		repo:      repo,
		encrypter: encrypter,
		prod:      prod,
	}
}
