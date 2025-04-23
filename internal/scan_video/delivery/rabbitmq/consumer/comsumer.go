package consumer

import (
	"context"
	"encoding/json"

	rmqDelivery "github.com/nguyenthanhtrung001/reup/internal/scan_video/delivery/rabbitmq"
	"github.com/nguyenthanhtrung001/reup/internal/scan_video/usecase"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (c Consumer) Consume() {
	go c.consume(rmqDelivery.ScanVideoOldExc, rmqDelivery.ScanVideoOldQueue, c.ScanVideoOld)
	go c.consume(rmqDelivery.ScanVideoNewExc, rmqDelivery.ScanVideoNewQueue, c.ScanVideoNew)
	go c.consume(rmqDelivery.ScanVideoManualExc, rmqDelivery.ScanVideoManualQueue, c.ScanVideoManual)

}

func (c Consumer) ScanVideoNew(d amqp.Delivery) {
	ctx := context.Background()

	var msg rmqDelivery.ScanDouyinVideosMsg
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		c.l.Error(ctx, "ScanVideo.delivery.rabbitmq.consumer.GetRandomBook.json.Unmarshal: %v", err)
		d.Ack(false)
		return
	}

	input := usecase.ScanDouyinVideosInput{
		ArrChannel: convertToUsecaseArrChannel(msg.ArrChannel),
		IsScanFull: msg.IsScanFull,
		Group:      msg.Group,
	}
	c.uc.SentScanDouyinVideos(input.ArrChannel, input.IsScanFull)

	d.Ack(false)
}
func (c Consumer) ScanVideoOld(d amqp.Delivery) {
	ctx := context.Background()

	var msg rmqDelivery.ScanDouyinVideosMsg
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		c.l.Error(ctx, "ScanVideo.delivery.rabbitmq.consumer.GetRandomBook.json.Unmarshal: %v", err)
		d.Ack(false)
		return
	}
	input := usecase.ScanDouyinVideosInput{
		ArrChannel: convertToUsecaseArrChannel(msg.ArrChannel),
		IsScanFull: msg.IsScanFull,
		Group:      msg.Group,
	}
	c.uc.SentScanDouyinVideos(input.ArrChannel, input.IsScanFull)

	d.Ack(false)
}
func (c Consumer) ScanVideoManual(d amqp.Delivery) {
	ctx := context.Background()

	var msg rmqDelivery.ScanDouyinVideosMsg
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		c.l.Error(ctx, "ScanVideo.delivery.rabbitmq.consumer.GetRandomBook.json.Unmarshal: %v", err)
		d.Ack(false)
		return
	}
	input := usecase.ScanDouyinVideosInput{
		ArrChannel: convertToUsecaseArrChannel(msg.ArrChannel),
		IsScanFull: msg.IsScanFull,
		Group:      msg.Group,
	}

	c.uc.SentScanDouyinVideos(input.ArrChannel, input.IsScanFull)

	d.Ack(false)
}

// convertToUsecaseArrChannel converts a slice of rmqDelivery.ArrChannel to a slice of usecase.ArrChannel
func convertToUsecaseArrChannel(channels []rmqDelivery.ArrChannel) []usecase.ArrChannel {
	result := make([]usecase.ArrChannel, len(channels))
	for i, ch := range channels {
		result[i] = usecase.ArrChannel{
			SpaceId:   ch.SpaceId,
			ChannelId: ch.ChannelId,
		}
	}
	return result
}
