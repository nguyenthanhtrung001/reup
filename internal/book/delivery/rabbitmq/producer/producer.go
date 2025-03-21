package producer

import (
	rabb "book-store/internal/book/delivery/rabbitmq"
	"book-store/pkg/rabbitmq"
	"context"
	"encoding/json"
)

// PubBookRandom implements Producer.
func (p implProducer) PubBookRandom(ctx context.Context, msg rabb.BookMsg) error {
	body, err := json.Marshal(msg)
	if err != nil {
		p.l.Errorf(ctx, "control.delivery.rabbitmq.producer.PubResetUnCompletedTask.json.Marshal: %v", err)
		return err
	}

	return p.bookRandomWriter.Publish(ctx, rabbitmq.PublishArgs{
		Exchange: rabb.GetRandomBookExcName,
		Msg: rabbitmq.Publishing{
			Body:        body,
			ContentType: rabbitmq.ContentTypePlainText,
		},
	})
}
