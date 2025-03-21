package schedule

import (
	"book-store/pkg/cron"
	pkgCrt "book-store/pkg/encrypter"
	pkgLog "book-store/pkg/log"
	"book-store/pkg/mongo"
	"book-store/pkg/rabbitmq"
)

type Schedule struct {
	cron      cron.Cron
	l         pkgLog.Logger
	db        mongo.Database
	encrypter pkgCrt.Encrypter
	conn      *rabbitmq.Connection
}
type Config struct {
	Database  mongo.Database
	Encrypter pkgCrt.Encrypter
	AMQPConn  *rabbitmq.Connection
}

func New(l pkgLog.Logger, cfg Config) Schedule {
	return Schedule{
		cron:      cron.New(),
		l:         l,
		db:        cfg.Database,
		encrypter: cfg.Encrypter,
		conn:      cfg.AMQPConn,
	}
}
