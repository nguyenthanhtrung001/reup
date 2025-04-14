package consumer

import (
	"context"
	"encoding/json"

	rmqDelivery "reup/internal/scan_video/delivery/rabbitmq"
	"reup/internal/scan_video/usecase"

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
	c.l.Info(ctx, "QUEUE NEW:", msg.SecUserID)

	input := usecase.ScanDouyinVideosInput{
		Mid:        msg.Mid,
		SecUserID:  msg.SecUserID,
		VideoCount: msg.VideoCount,
		NewFlag:    msg.NewFlag,
		Group:      msg.Group,
		DomainAPI:  msg.DomainAPI,
	}
	c.uc.ScanDouyinVideos(ctx, input)

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
	c.l.Info(ctx, "QUEUE OLD:", msg.SecUserID)

	input := usecase.ScanDouyinVideosInput{
		Mid:        msg.Mid,
		SecUserID:  msg.SecUserID,
		VideoCount: msg.VideoCount,
		NewFlag:    msg.NewFlag,
		Group:      msg.Group,
		DomainAPI:  msg.DomainAPI,
	}
	c.uc.ScanDouyinVideos(ctx, input)

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
	c.l.Info(ctx, "QUEUE MANUAL:", msg.SecUserID)

	input := usecase.ScanDouyinVideosInput{
		Mid:        msg.Mid,
		SecUserID:  msg.SecUserID,
		VideoCount: msg.VideoCount,
		NewFlag:    msg.NewFlag,
		Group:      msg.Group,
		DomainAPI:  msg.DomainAPI,
	}
	c.uc.ScanDouyinVideos(ctx, input)

	d.Ack(false)
}
