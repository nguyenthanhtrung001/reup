package usecase

import (
	"book-store/internal/book/repository"
	"book-store/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc implUseCase) CreateBook(ctx context.Context, input BookInput) error {
	newBook := repository.CreateOption{
		Title:       input.Title,
		Author:      input.Author,
		PublishedAt: input.PublishedAt,
	}
	return uc.repo.CreateBook(ctx, newBook)

}
func isValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

// DeleteBook implements UseCase.
func (uc implUseCase) DeleteBook(ctx context.Context, idBook string) error {
	_, err := uc.GetBookBy(ctx, idBook)
	if err != nil {
		return err
	}
	return uc.repo.DeleteBook(ctx, idBook)
}

// GetBookBy implements UseCase.
func (uc implUseCase) GetBookBy(ctx context.Context, id string) (BookOutput, error) {

	if !isValidObjectID(id) {
		return BookOutput{}, ErrInvalidID
	}

	book, err := uc.repo.DetailBook(ctx, id)
	if err != nil {
		uc.l.Errorf(ctx, "uc.GetBookBy - Error fetching book: %v", err)
		return BookOutput{}, ErrNotFound
	}

	return BookOutput{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		PublishedAt: book.PublishedAt,
		CreateAt:    book.CreateAt,
		UpdateAt:    book.UpdateAt,
	}, nil
}

// UpdateBook implements UseCase.
func (uc implUseCase) UpdateBook(ctx context.Context, i UpdateInput) error {
	_, err := uc.GetBookBy(ctx, i.ID)
	if err != nil {
		return err
	}
	updateBook := repository.UpdateOpition{
		ID:          i.ID,
		Title:       i.Title,
		Author:      i.Author,
		PublishedAt: i.PublishedAt,
	}
	return uc.repo.UpdateBook(ctx, updateBook)

}

// List implements UseCase.
func (uc implUseCase) List(ctx context.Context, sc models.Scope, i ListBookInput) (ListBookOutput, error) {
	op := repository.ListOptions{
		PaginatorQuery: i.PaginatorQuery,
		Filter: repository.ListFilterOptions{
			Title:  i.Filter.Title,
			Author: i.Filter.Author,
		},
	}
	books, paginator, err := uc.repo.ListBook(ctx, sc, op)
	if err != nil {
		uc.l.Errorf(ctx, "book.usecase.List.repo.List: %v", err)
		return ListBookOutput{}, err
	}
	return ListBookOutput{
		Books:     books,
		Paginator: paginator,
	}, nil
}
