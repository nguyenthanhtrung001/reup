package job

import (
	"github.com/nguyenthanhtrung001/reup/internal/scan_video/usecase"
	"github.com/nguyenthanhtrung001/reup/pkg/cron"
	pkgLog "github.com/nguyenthanhtrung001/reup/pkg/log"
)

type Handler struct {
	l    pkgLog.Logger
	uc   usecase.UseCase
	cron cron.Cron
}

func New(l pkgLog.Logger, uc usecase.UseCase, cronJ cron.Cron) Handler {
	return Handler{
		l:    l,
		cron: cronJ,
		uc:   uc,
	}
}
