package job

import (
	"reup/internal/book/usecase"
	"reup/pkg/cron"
	pkgLog "reup/pkg/log"
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
