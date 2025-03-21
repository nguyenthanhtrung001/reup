package repository

import (
	"book-store/internal/models"
	"book-store/pkg/paginator"
	"context"
)

type Repository interface {
	Book
}

type Book interface {
	CreateBook(ctx context.Context, input CreateOption) error
	UpdateBook(ctx context.Context, input UpdateOpition) error
	DetailBook(ctx context.Context, id string) (DetailOutputOption, error)
	DeleteBook(ctx context.Context, id string) error
	ListBook(ctx context.Context, sc models.Scope, opt ListOptions) ([]models.Book, paginator.Paginator, error)
	GetRandomBook() (models.Book, error)
}
