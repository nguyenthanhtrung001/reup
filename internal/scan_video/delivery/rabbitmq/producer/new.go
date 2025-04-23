package producer

import (
	"context"

	rabb "github.com/nguyenthanhtrung001/reup/internal/scan_video/delivery/rabbitmq"
	pkgLog "github.com/nguyenthanhtrung001/reup/pkg/log"
	"github.com/nguyenthanhtrung001/reup/pkg/rabbitmq"
)

type Producer interface {
	PubScanVideoOld(ctx context.Context, msg rabb.ScanDouyinVideosMsg) error
	PubScanVideoNew(ctx context.Context, msg rabb.ScanDouyinVideosMsg) error
	PubScanVideoManual(ctx context.Context, msg rabb.ScanDouyinVideosMsg) error
	Run() error
	Close()
}

type implProducer struct {
	l                     pkgLog.Logger
	conn                  *rabbitmq.Connection
	scanvideoOldWriter    *rabbitmq.Channel
	scanvideoNewdWriter   *rabbitmq.Channel
	scanvideoManualWriter *rabbitmq.Channel
}

func New(l pkgLog.Logger, conn *rabbitmq.Connection) Producer {
	return &implProducer{
		l:    l,
		conn: conn,
	}
}
