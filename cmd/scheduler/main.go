package main

import (
	"github.com/nguyenthanhtrung001/reup/config"
	"github.com/nguyenthanhtrung001/reup/internal/appconfig/mongo"
	"github.com/nguyenthanhtrung001/reup/internal/schedule"

	"context"

	pkgCrt "github.com/nguyenthanhtrung001/reup/pkg/encrypter"
	pkgLog "github.com/nguyenthanhtrung001/reup/pkg/log"
	"github.com/nguyenthanhtrung001/reup/pkg/rabbitmq"
)

func main() {

	ctx := context.Background()
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file sch: %v", err)
	// }
	// Load config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	l := pkgLog.InitializeZapLogger(pkgLog.ZapConfig{
		Level:    cfg.Logger.Level,
		Mode:     cfg.Logger.Mode,
		Encoding: cfg.Logger.Encoding,
	})
	amqpConn, err := rabbitmq.Dial(cfg.RabbitMQConfig.URL, true)
	if err != nil {

		panic(err)
	}
	defer amqpConn.Close()
	crp := pkgCrt.NewEncrypter(cfg.Encrypter.Key)

	client, err := mongo.Connect(cfg.Mongo, crp)
	if err != nil {
		l.Fatalf(ctx, "Failed to connect to MongoDB: %v", err)
	}
	defer mongo.Disconnect(client)

	db := client.Database(cfg.Mongo.Database)

	err = schedule.New(l, schedule.Config{
		Database:  db,
		Encrypter: crp,
		AMQPConn:  amqpConn,
		Telegram: schedule.TeleCredentials{
			BotKey: cfg.Telegram.BotKey,
			ChatIDs: schedule.ChatIDs{
				ReportBug:     cfg.Telegram.ChatIDs.ReportBug,
				ReportPayment: cfg.Telegram.ChatIDs.ReportPayment,
			},
		},
	}).Start()

	if err != nil {
		l.Fatalf(ctx, "Failed to create scheduler: %v", err)
		panic(err)
	}
}
