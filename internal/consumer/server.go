package consumer

import (
	"context"

	"reup/pkg/jwt"
	"reup/pkg/log"
	"reup/pkg/mongo"
	"reup/pkg/rabbitmq"
	"reup/pkg/redis"
	"reup/pkg/telegram"

	bookConsumer "reup/internal/book/delivery/rabbitmq/consumer"
	bookProd "reup/internal/book/delivery/rabbitmq/producer"
	bookMongo "reup/internal/book/repository/mongo"
	bookUseCase "reup/internal/book/usecase"

	pkgCrt "reup/pkg/encrypter"
)

// Server is the consumer server
type Server struct {
	l         log.Logger
	conn      *rabbitmq.Connection
	db        mongo.Database
	redis     redis.Client
	encrypter pkgCrt.Encrypter
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

// NewServer creates a new consumer server
func NewServer(l log.Logger,
	conn *rabbitmq.Connection,
	db mongo.Database,
	redis redis.Client,
	encrypter pkgCrt.Encrypter,
	telegram TeleCredentials,
) Server {
	return Server{
		l:         l,
		conn:      conn,
		db:        db,
		redis:     redis,
		encrypter: encrypter,
		telegram:  telegram,
	}
}

// Run runs the consumer server
func (s Server) Run() error {

	chatIDs := telegram.ChatIDs{
		ReportBug:     s.telegram.ChatIDs.ReportBug,
		ReportPayment: s.telegram.ChatIDs.ReportPayment,
	}
	telegram := telegram.New(s.telegram.BotKey, chatIDs)
	// Producer
	bookProd := bookProd.New(s.l, s.conn)
	if err := bookProd.Run(); err != nil {
		s.l.Fatal(context.Background(), err)
		return err
	}
	// Cac usecase su dung tele them o day
	s.l.Fatal(context.Background(), telegram)
	// Repositories

	bookMongo := bookMongo.New(s.l, s.db, jwt.JWTMaker{})

	// UseCases

	bookUseCase := bookUseCase.New(s.l, bookMongo, pkgCrt.NewEncrypter(""), bookProd)

	// Consumer
	var forever chan bool
	bookConsumer.NewConsumer(s.l, s.conn, bookUseCase).Consume()

	<-forever

	return nil
}
