package producer

import (
	rmqDelivery "reup/internal/scan_video/delivery/rabbitmq"
	rmqPkg "reup/pkg/rabbitmq"
)

// Run runs the producer
func (p *implProducer) Run() (err error) {

	if p.scanvideoOldWriter, err = p.getWriter(rmqDelivery.ScanVideoOldExc); err != nil {
		return
	}
	if p.scanvideoManualWriter, err = p.getWriter(rmqDelivery.ScanVideoManualExc); err != nil {

		return
	}
	if p.scanvideoNewdWriter, err = p.getWriter(rmqDelivery.ScanVideoNewExc); err != nil {
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
