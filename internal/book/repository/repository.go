package repository

import (
	"context"
	"reup/internal/models"
	"reup/pkg/paginator"
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
