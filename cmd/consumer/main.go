package main

import (
	"context"

	"book-store/config"
	"book-store/internal/appconfig/mongo"
	"book-store/internal/appconfig/redis"
	"book-store/internal/consumer"
	pkgCrt "book-store/pkg/encrypter"
	pkgLog "book-store/pkg/log"
	"book-store/pkg/rabbitmq"
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

	crp := pkgCrt.NewEncrypter(cfg.Encrypter.Key)

	client, err := mongo.Connect(cfg.Mongo, crp)
	if err != nil {
		l.Fatalf(ctx, "Failed to connect to MongoDB: %v", err)
	}
	defer mongo.Disconnect(client)

	db := client.Database(cfg.Mongo.Database)

	conn, err := rabbitmq.Dial(cfg.RabbitMQConfig.URL, true)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	redisClient, err := redis.Connect(cfg.RedisConfig)
	if err != nil {
		panic(err)
	}

	if err := consumer.NewServer(l, conn, db, redisClient, crp).Run(); err != nil {
		l.Fatalf(ctx, "Failed to run consumer server: %v", err)
	}
	l.Info(ctx, "Consumer server is running successfully")
}
