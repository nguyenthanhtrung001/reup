package consumer

import (
	"context"

	"reup/pkg/jwt"
	"reup/pkg/log"
	"reup/pkg/mongo"
	"reup/pkg/rabbitmq"
	"reup/pkg/redis"
	"reup/pkg/telegram"

	scanVideoConsumer "reup/internal/scan_video/delivery/rabbitmq/consumer"
	scanVideoProd "reup/internal/scan_video/delivery/rabbitmq/producer"
	scanVideoMongo "reup/internal/scan_video/repository/mongo"
	scanVideoUseCase "reup/internal/scan_video/usecase"

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
	GroupChat1    int64
	GroupChat2    int64
	GroupChat3    int64
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
	scanVideoProd := scanVideoProd.New(s.l, s.conn)
	if err := scanVideoProd.Run(); err != nil {
		s.l.Fatal(context.Background(), err)
		return err
	}
	// Cac usecase su dung tele them o day
	//s.l.Fatal(context.Background(), telegram)
	s.l.Info(context.Background(), telegram)
	// Repositories

	scanVideoMongo := scanVideoMongo.New(s.l, s.db, jwt.JWTMaker{})

	// UseCases

	scanVideoUseCase := scanVideoUseCase.New(s.l, scanVideoMongo, pkgCrt.NewEncrypter(""), scanVideoProd, telegram, scanVideoUseCase.TeleChat{
		NotifiChatID: s.telegram.ReportPayment,
		GroupChat1:   s.telegram.GroupChat1,
		GroupChat2:   s.telegram.GroupChat2,
		GroupChat3:   s.telegram.GroupChat3,
	})

	// Consumer
	var forever = make(chan bool)
	scanVideoConsumer.NewConsumer(s.l, s.conn, scanVideoUseCase).Consume()

	<-forever

	return nil
}
