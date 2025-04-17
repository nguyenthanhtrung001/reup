package http

import (
	"reup/internal/middleware"

	"github.com/gin-gonic/gin"
)

func MapScanVideoRoutes(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	// r.Use(mw.Auth()).Use(mw.AuthReseller(true))
	r.POST("/test", h.HandleDouyinWebhook)

	r.GET("/proxy_scan/done_all", h.DoneAllProxyScan)
	r.GET("/proxy_scan/done", h.DoneProxyScan)
	r.GET("/proxy_scan/all", h.GetAllProxyScan)
	r.GET("/proxy_scan/random", h.GetProxyScanRandom)

}
