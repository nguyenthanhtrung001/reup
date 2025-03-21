package http

import (
	"reup/internal/book/usecase"
	"reup/internal/models"
	"reup/pkg/paginator"
	"reup/pkg/response"

	"github.com/gin-gonic/gin"
)

// AddBook implements Handler.
func (h handler) AddBook(c *gin.Context) {
	ctx := c.Request.Context()
	_, err := h.authenticateUser(c)
	if err != nil {
		h.l.Errorf(ctx, "book.http.addbook.auth:%v", err)
		response.Error(c, err)
		return
	}
	rqBody, err := h.processBookReq(c)
	if err != nil {
		h.l.Errorf(ctx, "book.http.addbook.processBookReq: %v", err)
		response.Error(c, err)
		return
	}
	newBook := usecase.BookInput{
		Title:       rqBody.Title,
		Author:      rqBody.Author,
		PublishedAt: rqBody.PublishedAt,
	}
	if err := h.uc.CreateBook(ctx, newBook); err != nil {
		h.l.Errorf(ctx, "book.http.addbook.uc")
		response.Error(c, err)
	}
	response.OK(c, newBook)

}

// DeleteBook implements Handler.
func (h handler) DeleteBook(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	err := h.uc.DeleteBook(ctx, id)
	if err != nil {
		h.l.Errorf(ctx, "book.http.uc.delete")
		response.Error(c, err)
		return
	}
	response.OK(c, "Delete successfull")
}

// GetBook implements Handler.
func (h handler) GetBook(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	book, err := h.uc.GetBookBy(ctx, id)
	if err != nil {
		h.l.Errorf(ctx, "book.http.uc.getDetail")
		response.Error(c, err)
		return
	}
	response.OK(c, book)
}

// UpdateBook implements Handler.
func (h handler) UpdateBook(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	rqBody, err := h.processBookReq(c)
	if err != nil {
		h.l.Errorf(ctx, "book.http.updatebook.processBookReq: %v", err)
		response.Error(c, err)
		return
	}
	updateBook := usecase.UpdateInput{
		ID:          id,
		Title:       rqBody.Title,
		Author:      rqBody.Author,
		PublishedAt: rqBody.PublishedAt,
	}
	err = h.uc.UpdateBook(ctx, updateBook)
	if err != nil {
		h.l.Errorf(ctx, "book.http.updatebook.uc: %v", err)
		response.Error(c, err)
		return
	}
	book, err := h.uc.GetBookBy(ctx, id)
	if err != nil {
		h.l.Errorf(ctx, "book.http.updatebook.uc.GetBookBy: %v", err)
		response.Error(c, err)
		return
	}
	response.OK(c, book)

}

// ListBook implements Handler.
func (h handler) ListBook(c *gin.Context) {
	ctx := c.Request.Context()
	req, err := h.processListBookReq(c)
	if err != nil {
		h.l.Errorf(ctx, "book.http.List.processListBookReq: %v", err)
		response.Error(c, err)
		return
	}
	var pagQuery paginator.PaginatorQuery
	if err := c.ShouldBindQuery(&pagQuery); err != nil {
		h.l.Errorf(ctx, "user.http.Get.gin.ShouldBindQuery: %v", err)
		response.Error(c, errWrongPaginationQuery)
		return
	}
	pagQuery.Adjust()
	books, err := h.uc.List(ctx, models.Scope{}, usecase.ListBookInput{
		PaginatorQuery: pagQuery,
		Filter: usecase.ListFilterOptions{
			Title:  req.Title,
			Author: req.Author,
		},
	})
	if err != nil {
		h.l.Errorf(ctx, "book.http.List.userUC.List: %v", err)
		response.Error(c, err)
		return
	}
	response.OK(c, newListBookResponse(books))
}
