package config

import "github.com/caarlos0/env/v9"

type Config struct {
	HTTPServer     HTTPServerConfig
	Logger         LoggerConfig
	JWT            JWTConfig
	Mongo          MongoConfig
	Encrypter      EncrypterConfig
	RabbitMQConfig RabbitMQConfig
	RedisConfig    RedisConfig
	Telegram       TelegramConfig
}

type RedisConfig struct {
	RedisAddr     string `env:"REDIS_HOST"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       string `env:"REDIS_DATABASE"`
	MinIdleConns  int    `env:"REDIS_MIN_IDLE_CONNS"`
	PoolSize      int    `env:"REDIS_POOL_SIZE"`
	PoolTimeout   int    `env:"REDIS_POOL_TIMEOUT"`
	Password      string `env:"REDIS_PASSWORD"`
	DB            int    `env:"REDIS_DATABASE"`
}

type RabbitMQConfig struct {
	URL string `env:"RABBITMQ_URL"`
}

type JWTConfig struct {
	SecretKey string `env:"JWT_SECRET"`
}

type HTTPServerConfig struct {
	Port int    `env:"API_PORT" envDefault:"8889"`
	Mode string `env:"API_MODE" envDefault:"debug"`
}

type LoggerConfig struct {
	Level    string `env:"LOGGER_LEVEL" envDefault:"debug"`
	Mode     string `env:"LOGGER_MODE" envDefault:"development"`
	Encoding string `env:"LOGGER_ENCODING" envDefault:"console"`
}

type MongoConfig struct {
	Database       string `env:"MONGODB_DATABASE"`
	ENCODED_URI    string `env:"MONGODB_ENCODED_URI"`
	ENABLE_MONITOR bool   `env:"MONGODB_ENABLE_MONITORING" envDefault:"false"`
}

type EncrypterConfig struct {
	Key string `env:"ENCRYPT_KEY"`
}

type TelegramConfig struct {
	BotKey  string `env:"TELEGRAM_BOT_KEY"`
	ChatIDs TeleChatIDs
}
type TeleChatIDs struct {
	ReportBug     int64 `env:"TELEGRAM_REPORT_BUG"`
	ReportPayment int64 `env:"TELEGRAM_REPORT_PAYMENT"`
	GroupChat1    int64 `env:"TELEGRAM_GROUP_1"`
	GroupChat2    int64 `env:"TELEGRAM_GROUP_2"`
	GroupChat3    int64 `env:"TELEGRAM_GROUP_3"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
