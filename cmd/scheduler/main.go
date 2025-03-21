package main

import (
	"book-store/config"
	"book-store/internal/appconfig/mongo"
	"book-store/internal/schedule"
	pkgCrt "book-store/pkg/encrypter"
	pkgLog "book-store/pkg/log"
	"book-store/pkg/rabbitmq"
	"context"
)

func main() {

	ctx := context.Background()

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
	}).Start()

	if err != nil {
		l.Fatalf(ctx, "Failed to create scheduler: %v", err)
		panic(err)
	}
}
