package http

import (
	"reup/internal/middleware"

	"github.com/gin-gonic/gin"
)

func MapUserRoutes(group *gin.RouterGroup, handler Handler, mw middleware.Middleware) {
	group.POST("/register", handler.Register)
	group.POST("/login", handler.Login)
	group.Use(mw.Auth())
	group.GET("/profile", handler.Profile)
	group.PUT("/profile", handler.UpdateProfile)
	group.PUT("/password", handler.UpdatePassword)
}
