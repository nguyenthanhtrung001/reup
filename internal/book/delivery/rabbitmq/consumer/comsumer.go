package consumer

import (
	"context"
	"encoding/json"

	rmqDelivery "reup/internal/book/delivery/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (c Consumer) Consume() {
	go c.consume(rmqDelivery.RandomBookExc, rmqDelivery.GetRandomBookQueue, c.getRandomBook)

}

func (c Consumer) getRandomBook(d amqp.Delivery) {
	ctx := context.Background()

	var msg rmqDelivery.BookMsg
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		c.l.Error(ctx, "book.delivery.rabbitmq.consumer.GetRandomBook.json.Unmarshal: %v", err)
		d.Ack(false)
		return
	}
	c.l.Info(ctx, "du lieu gui vao:", msg.Msg)
	// var book models.Book
	// book, err = c.uc.GetRandomBookTask()
	// if err != nil {
	// 	c.l.Warnf(ctx, "books.delivery.rabbitmq.consumer.resetTotalTaGetRandomBookskWorker.uc.GetRandomBook: %v", err)
	// 	d.Ack(false)
	// 	return
	// }
	// c.l.Info(ctx, "kết quả:", msg)
	// c.l.Info(ctx, book.Title)

	d.Ack(false)
}
