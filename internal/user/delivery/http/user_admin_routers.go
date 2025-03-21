package http

import (
	"book-store/internal/middleware"

	"github.com/gin-gonic/gin"
)

func MapUserAdminRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.Use(mw.Auth()).Use(mw.AuthReseller(true))
	r.GET("/list", h.List)
	r.GET("/:id", h.DetailAdmin)
	r.PUT("/:id", h.UpdateByAdmin)
	r.POST("/:id/api-key", h.CreateApiKey)
}
