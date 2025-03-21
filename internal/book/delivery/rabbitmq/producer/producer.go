package producer

import (
	"context"
	"encoding/json"
	rabb "reup/internal/book/delivery/rabbitmq"
	"reup/pkg/rabbitmq"
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
