package http

import (
	"book-store/internal/middleware"

	"github.com/gin-gonic/gin"
)

func MapbookRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	// r.Use(mw.Auth()).Use(mw.AuthReseller(true))
	r.Use(mw.Auth()).Use(mw.AuthUser())
	r.GET("/test", h.test)
	r.GET("/detail/:id", h.GetBook)
	r.GET("/list", h.ListBook)
	r.POST("/add-book", h.AddBook)
	r.PATCH("/update-book/:id", h.UpdateBook)
	r.DELETE("/delete-book/:id", h.DeleteBook)
}
