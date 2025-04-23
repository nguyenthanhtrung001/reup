package http

import (
	"github.com/nguyenthanhtrung001/reup/internal/scan_video/usecase"
	"github.com/nguyenthanhtrung001/reup/pkg/log"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	ScanVideoHandler
	ProxySxanHandler
}

type ScanVideoHandler interface {
	HandleDouyinWebhook(c *gin.Context)
}
type ProxySxanHandler interface {
	DoneAllProxyScan(c *gin.Context)
	DoneProxyScan(c *gin.Context)
	GetAllProxyScan(c *gin.Context)
	GetProxyScanRandom(c *gin.Context)
	InsertProxyScan(c *gin.Context)

	GetQuestHandler(c *gin.Context)
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
