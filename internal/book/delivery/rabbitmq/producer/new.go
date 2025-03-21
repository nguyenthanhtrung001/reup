package producer

import (
	"context"
	rabb "reup/internal/book/delivery/rabbitmq"
	pkgLog "reup/pkg/log"
	"reup/pkg/rabbitmq"
)

type Producer interface {
	PubBookRandom(ctx context.Context, msg rabb.BookMsg) error
	Run() error
	Close()
}

type implProducer struct {
	l                pkgLog.Logger
	conn             *rabbitmq.Connection
	bookRandomWriter *rabbitmq.Channel
}

func New(l pkgLog.Logger, conn *rabbitmq.Connection) Producer {
	return &implProducer{
		l:    l,
		conn: conn,
	}
}
