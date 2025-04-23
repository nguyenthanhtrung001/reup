package httpserver

import (
	pkgCrt "github.com/nguyenthanhtrung001/reup/pkg/encrypter"
	pkgLog "github.com/nguyenthanhtrung001/reup/pkg/log"
	"github.com/nguyenthanhtrung001/reup/pkg/mongo"
	"github.com/nguyenthanhtrung001/reup/pkg/rabbitmq"

	// "github.com/nguyenthanhtrung001/reup/pkg/rabbitmq"
	"github.com/nguyenthanhtrung001/reup/pkg/redis"

	"github.com/gin-gonic/gin"
)

const productionMode = "production"

var ginMode = gin.DebugMode

type HTTPServer struct {
	gin          *gin.Engine
	l            pkgLog.Logger
	port         int
	database     mongo.Database
	jwtSecretKey string
	mode         string
	amqpConn     *rabbitmq.Connection
	redis        redis.Client

	telegram  TeleCredentials
	encrypter pkgCrt.Encrypter
	secretKey string
}
type Config struct {
	Port         int
	JWTSecretKey string
	Database     mongo.Database
	Mode         string
	AMQPConn     *rabbitmq.Connection
	Redis        redis.Client
	Telegram     TeleCredentials
	Encrypter    pkgCrt.Encrypter
	SecretKey    string
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

func New(l pkgLog.Logger, cfg Config) *HTTPServer {
	if cfg.Mode == productionMode {
		ginMode = gin.ReleaseMode
	}

	gin.SetMode(ginMode)

	return &HTTPServer{
		l:            l,
		gin:          gin.Default(),
		port:         cfg.Port,
		database:     cfg.Database,
		jwtSecretKey: cfg.JWTSecretKey,
		mode:         cfg.Mode,
		amqpConn:     cfg.AMQPConn,
		redis:        cfg.Redis,
		telegram:     cfg.Telegram,
		encrypter:    cfg.Encrypter,
		secretKey:    cfg.SecretKey,
	}
}
