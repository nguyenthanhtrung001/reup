package http

import (
	"book-store/internal/book/usecase"
	"book-store/pkg/log"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	EncryHandler
	BookHandler
}

type EncryHandler interface {
	test(c *gin.Context)
}

type BookHandler interface {
	AddBook(c *gin.Context)
	GetBook(c *gin.Context)
	UpdateBook(c *gin.Context)
	DeleteBook(c *gin.Context)
	ListBook(c *gin.Context)
}

type handler struct {
	l  log.Logger
	uc usecase.UseCase
}

func New(l log.Logger, uc usecase.UseCase) Handler {
	return handler{
		l:  l,
		uc: uc,
	}

}
