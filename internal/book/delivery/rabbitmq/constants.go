package rabbitmq

import "book-store/pkg/rabbitmq"

const (
	GetRandomBookExcName = "book_random_exc"
	GetRandomBookQueue   = "book_random_queue"
)

var (
	RandomBookExc = rabbitmq.ExchangeArgs{
		Name:       GetRandomBookExcName,
		Type:       rabbitmq.ExchangeTypeFanout,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	}
)
