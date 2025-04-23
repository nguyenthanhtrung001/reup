package consumer

import (
	"github.com/nguyenthanhtrung001/reup/internal/scan_video/usecase"
	pkgLog "github.com/nguyenthanhtrung001/reup/pkg/log"
	"github.com/nguyenthanhtrung001/reup/pkg/rabbitmq"
)

// Consumer represents a consumer

type Consumer struct {
	l    pkgLog.Logger
	conn *rabbitmq.Connection
	uc   usecase.UseCase
}

// NewConsumer creates a new consumer
func NewConsumer(l pkgLog.Logger, conn *rabbitmq.Connection, uc usecase.UseCase) Consumer {
	return Consumer{
		l:    l,
		conn: conn,
		uc:   uc,
	}
}
