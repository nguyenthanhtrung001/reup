package main

import (
	"context"

	"github.com/nguyenthanhtrung001/reup/config"
	"github.com/nguyenthanhtrung001/reup/internal/appconfig/mongo"
	"github.com/nguyenthanhtrung001/reup/internal/appconfig/redis"
	"github.com/nguyenthanhtrung001/reup/internal/consumer"
	pkgCrt "github.com/nguyenthanhtrung001/reup/pkg/encrypter"
	pkgLog "github.com/nguyenthanhtrung001/reup/pkg/log"
	"github.com/nguyenthanhtrung001/reup/pkg/rabbitmq"
)

func main() {

	ctx := context.Background()
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file con: %v", err)
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

	if err := consumer.NewServer(l, conn, db, redisClient, crp,
		consumer.TeleCredentials{
			BotKey: cfg.Telegram.BotKey,
			ChatIDs: consumer.ChatIDs{
				ReportBug:     cfg.Telegram.ChatIDs.ReportBug,
				ReportPayment: cfg.Telegram.ChatIDs.ReportPayment,
				GroupChat1:    cfg.Telegram.ChatIDs.GroupChat1,
				GroupChat2:    cfg.Telegram.ChatIDs.GroupChat2,
				GroupChat3:    cfg.Telegram.ChatIDs.GroupChat3,
			}}).Run(); err != nil {
		l.Fatalf(ctx, "Failed to run consumer server: %v", err)
	}
	l.Info(ctx, "Consumer server is running successfully")
}
