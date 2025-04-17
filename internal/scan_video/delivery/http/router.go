package http

import (
	"reup/internal/middleware"

	"github.com/gin-gonic/gin"
)

func MapScanVideoRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	// r.Use(mw.Auth()).Use(mw.AuthReseller(true))
	r.POST("/test", h.HandleDouyinWebhook)

}
