package http

import (
	"reup/internal/book/usecase"
	"reup/internal/models"
	"reup/pkg/paginator"
)

type bookRequest struct {
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	PublishedAt string `json:"published_at"`
}
type listBookReq struct {
	Title  string `json:"title" form:"title"`
	Author string `json:"author" form:"author"`
}
type newBook struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishedAt string `json:"published_at"`
	CreateAt    string `json:"create_at"`
	UpdateAt    string `json:"update_at"`
}
type listBookResponse struct {
	Books []newBook        `json:"books"`
	Meta  listMetaResponse `json:"meta"`
}
type listMetaResponse struct {
	paginator.PaginatorResponse
}

func newBookResponse(book models.Book) newBook {
	return newBook{
		ID:          book.ID.Hex(),
		Author:      book.Author,
		Title:       book.Title,
		PublishedAt: book.PublishedAt,
		CreateAt:    book.CreateAt.Format("2006-01-02 15:04:05"),
		UpdateAt:    book.UpdateAt.Format("2006-01-02 15:04:05"),
	}
}
func newListBookResponse(books usecase.ListBookOutput) listBookResponse {
	items := make([]newBook, 0, len(books.Books))
	for _, v := range books.Books {
		items = append(items, newBookResponse(v))
	}
	return listBookResponse{
		Books: items,
		Meta: listMetaResponse{
			PaginatorResponse: books.Paginator.ToResponse(),
		},
	}
}
