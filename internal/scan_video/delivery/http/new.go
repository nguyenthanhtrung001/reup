package http

import (
	"reup/internal/scan_video/usecase"
	"reup/pkg/log"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	ScanVideoHandler
}

type ScanVideoHandler interface {
	HandleDouyinWebhook(c *gin.Context)
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
