package schedule

import (
	"github.com/nguyenthanhtrung001/reup/pkg/cron"
	pkgCrt "github.com/nguyenthanhtrung001/reup/pkg/encrypter"
	pkgLog "github.com/nguyenthanhtrung001/reup/pkg/log"
	"github.com/nguyenthanhtrung001/reup/pkg/mongo"
	"github.com/nguyenthanhtrung001/reup/pkg/rabbitmq"
)

type Scheduler struct {
	cron      cron.Cron
	l         pkgLog.Logger
	db        mongo.Database
	encrypter pkgCrt.Encrypter
	conn      *rabbitmq.Connection
	telegram  TeleCredentials
}
type TeleCredentials struct {
	BotKey string
	ChatIDs
}

type ChatIDs struct {
	ReportBug     int64
	ReportPayment int64
}

type Config struct {
	Database  mongo.Database
	Encrypter pkgCrt.Encrypter
	AMQPConn  *rabbitmq.Connection
	Telegram  TeleCredentials
}

func New(l pkgLog.Logger, cfg Config) Scheduler {
	return Scheduler{
		cron:      cron.New(),
		l:         l,
		db:        cfg.Database,
		encrypter: cfg.Encrypter,
		conn:      cfg.AMQPConn,
		telegram:  cfg.Telegram,
	}
}
