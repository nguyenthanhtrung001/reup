package main

import (
	"book-store/config"
	"book-store/internal/appconfig/mongo"
	"book-store/internal/appconfig/redis"
	"book-store/internal/httpserver"
	pkgCrt "book-store/pkg/encrypter"
	pkgLog "book-store/pkg/log"
	"book-store/pkg/rabbitmq"
)

// @title SSC Group API - SMM
// @description This is the API documentation for the SSC Group SMM service.
// @description Error codes:
// @description `4000 ("Invalid password"),`
// @description `4002 ("Invalid email"),`
// @description `4003 ("User already registered"),`
// @description `4004 ("User not found"),`
// @description `4005 ("Missing params, please check email, phone, fullname, password"),`
// @description `4006 ("Invalid phone"),`
// @description `6000 ("Wrong pagination query"),`
// @description `6001 ("Invalid body"),`
// @description `6002 ("Invalid validation"),`

// @version 1
// @host 10.10.10.186:8881

// @BasePath /
// @schemes http
func main() {

	// Load config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	crp := pkgCrt.NewEncrypter(cfg.Encrypter.Key)
	client, err := mongo.Connect(cfg.Mongo, crp)
	if err != nil {
		panic(err)
	}
	defer mongo.Disconnect(client)

	database := client.Database(cfg.Mongo.Database)

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

	redisClient, err := redis.Connect(cfg.RedisConfig)
	if err != nil {
		panic(err)
	}
	defer redisClient.Disconnect()

	srv := httpserver.New(l, httpserver.Config{
		Port:         cfg.HTTPServer.Port,
		Database:     database,
		JWTSecretKey: cfg.JWT.SecretKey,
		Mode:         cfg.HTTPServer.Mode,
		AMQPConn:     amqpConn,
		Redis:        redisClient,

		Telegram: httpserver.TeleCredentials{
			BotKey: cfg.Telegram.BotKey,
			ChatIDs: httpserver.ChatIDs{
				ReportBug:     cfg.Telegram.ChatIDs.ReportBug,
				ReportPayment: cfg.Telegram.ChatIDs.ReportPayment,
			},
		},
		Encrypter: crp,
		SecretKey: cfg.Encrypter.Key,
	})

	if err := srv.Run(); err != nil {
		panic(err)
	}
}
