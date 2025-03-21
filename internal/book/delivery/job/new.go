package job

import (
	"book-store/internal/book/usecase"
	"book-store/pkg/cron"
	pkgLog "book-store/pkg/log"
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
