package http

import (
	"reup/internal/middleware"

	"github.com/gin-gonic/gin"
)

func MapProxyRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	// r.Use(mw.Auth()).Use(mw.AuthReseller(true))
	r.POST("/done-all-proxy", h.DoneAllProxyScan)
	r.POST("/done-proxy", h.DoneProxyScan)
	r.GET("/proxy-all", h.GetAllProxyScan)
	r.GET("/proxy", h.GetProxyScanRandom)
	r.POST("/insert-proxy", h.InsertProxyScan)
}
