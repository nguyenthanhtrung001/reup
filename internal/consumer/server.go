package consumer

import (
	"context"

	"book-store/pkg/jwt"
	"book-store/pkg/log"
	"book-store/pkg/mongo"
	"book-store/pkg/rabbitmq"
	"book-store/pkg/redis"

	bookConsumer "book-store/internal/book/delivery/rabbitmq/consumer"
	bookProd "book-store/internal/book/delivery/rabbitmq/producer"
	bookMongo "book-store/internal/book/repository/mongo"
	bookUseCase "book-store/internal/book/usecase"

	pkgCrt "book-store/pkg/encrypter"
)

// Server is the consumer server
type Server struct {
	l         log.Logger
	conn      *rabbitmq.Connection
	db        mongo.Database
	redis     redis.Client
	encrypter pkgCrt.Encrypter
}

// NewServer creates a new consumer server
func NewServer(l log.Logger,
	conn *rabbitmq.Connection,
	db mongo.Database,
	redis redis.Client,
	encrypter pkgCrt.Encrypter,
) Server {
	return Server{
		l:         l,
		conn:      conn,
		db:        db,
		redis:     redis,
		encrypter: encrypter,
	}
}

// Run runs the consumer server
func (s Server) Run() error {

	// Producer
	bookProd := bookProd.New(s.l, s.conn)
	if err := bookProd.Run(); err != nil {
		s.l.Fatal(context.Background(), err)
		return err
	}

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
