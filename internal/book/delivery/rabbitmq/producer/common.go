package producer

import (
	rmqDelivery "reup/internal/book/delivery/rabbitmq"
	rmqPkg "reup/pkg/rabbitmq"
)

// Run runs the producer
func (p *implProducer) Run() (err error) {

	if p.bookRandomWriter, err = p.getWriter(rmqDelivery.RandomBookExc); err != nil {
		return
	}

	return
}

// Close closes the producer
func (p *implProducer) Close() {

}

func (p implProducer) getWriter(exchange rmqPkg.ExchangeArgs) (*rmqPkg.Channel, error) {
	ch, err := p.conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(exchange)
	if err != nil {
		return nil, err
	}

	return ch, nil
}
