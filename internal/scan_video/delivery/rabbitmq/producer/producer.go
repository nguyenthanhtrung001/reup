package producer

import (
	"context"
	"encoding/json"
	rabb "reup/internal/scan_video/delivery/rabbitmq"
	"reup/pkg/rabbitmq"
)

func (p implProducer) PubScanVideoNew(ctx context.Context, msg rabb.ScanDouyinVideosMsg) error {
	body, err := json.Marshal(msg)
	if err != nil {
		p.l.Errorf(ctx, "control.delivery.rabbitmq.producer.PubScanVideoNew.json.Marshal: %v", err)
		return err
	}

	return p.scanvideoNewdWriter.Publish(ctx, rabbitmq.PublishArgs{
		Exchange: rabb.ScanVideoNewExcName,
		Msg: rabbitmq.Publishing{
			Body:        body,
			ContentType: rabbitmq.ContentTypePlainText,
		},
	})
}
func (p implProducer) PubScanVideoManual(ctx context.Context, msg rabb.ScanDouyinVideosMsg) error {
	body, err := json.Marshal(msg)
	if err != nil {
		p.l.Errorf(ctx, "control.delivery.rabbitmq.producer.PubScanVideoManual.json.Marshal: %v", err)
		return err
	}

	return p.scanvideoManualWriter.Publish(ctx, rabbitmq.PublishArgs{
		Exchange: rabb.ScanVideoManualExcName,
		Msg: rabbitmq.Publishing{
			Body:        body,
			ContentType: rabbitmq.ContentTypePlainText,
		},
	})
}

func (p implProducer) PubScanVideoOld(ctx context.Context, msg rabb.ScanDouyinVideosMsg) error {
	body, err := json.Marshal(msg)
	if err != nil {
		p.l.Errorf(ctx, "control.delivery.rabbitmq.producer.PubScanVideoOld.json.Marshal: %v", err)
		return err
	}
	// p.l.Infof(ctx, "Sending message to RabbitMQ: %s", string(body))

	return p.scanvideoOldWriter.Publish(ctx, rabbitmq.PublishArgs{
		Exchange: rabb.ScanVideoOldExcName,
		Msg: rabbitmq.Publishing{
			Body:        body,
			ContentType: rabbitmq.ContentTypePlainText,
		},
	})
}
