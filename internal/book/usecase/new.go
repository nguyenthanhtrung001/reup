package usecase

import (
	prod "book-store/internal/book/delivery/rabbitmq/producer"
	"book-store/internal/book/repository"
	"book-store/internal/models"
	"book-store/pkg/encrypter"
	"book-store/pkg/log"
	"context"
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
