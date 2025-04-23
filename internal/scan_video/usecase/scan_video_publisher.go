package usecase

import (
	"context"

	"github.com/nguyenthanhtrung001/reup/internal/scan_video/delivery/rabbitmq"
)

func (uc implUseCase) pubScanVideoNewTask(ctx context.Context, msg ScanDouyinVideosInput) error {
	uc.l.Info(ctx, "========== đẩy queue video new ==========")
	err := uc.prod.PubScanVideoOld(ctx, rabbitmq.ScanDouyinVideosMsg{
		ArrChannel: convertToRabbitMQArrChannel(msg.ArrChannel),
		IsScanFull: msg.IsScanFull,
		Group:      msg.Group,
	})
	if err != nil {
		uc.l.Errorf(ctx, "book.usecase.pubScanVideoOldTask.prod.pubRandomBookTask: %v", err)
		return err
	}

	return nil
}
func (uc implUseCase) pubScanVideoManualTask(ctx context.Context, msg ScanDouyinVideosInput) error {
	uc.l.Info(ctx, "========== đẩy queue video manual ==========")
	err := uc.prod.PubScanVideoOld(ctx, rabbitmq.ScanDouyinVideosMsg{
		ArrChannel: convertToRabbitMQArrChannel(msg.ArrChannel),
		IsScanFull: msg.IsScanFull,
		Group:      msg.Group,
	})
	if err != nil {
		uc.l.Errorf(ctx, "book.usecase.pubScanVideoOldTask.prod.pubRandomBookTask: %v", err)
		return err
	}

	return nil
}

// convertToRabbitMQArrChannel converts []ArrChannel to []rabbitmq.ArrChannel.
func convertToRabbitMQArrChannel(channels []ArrChannel) []rabbitmq.ArrChannel {
	result := make([]rabbitmq.ArrChannel, len(channels))
	for i, ch := range channels {
		result[i] = rabbitmq.ArrChannel(ch)
	}
	return result
}

func (uc implUseCase) pubScanVideoOldTask(ctx context.Context, msg ScanDouyinVideosInput) error {
	uc.l.Info(ctx, "========== đẩy queue video old ==========")
	err := uc.prod.PubScanVideoOld(ctx, rabbitmq.ScanDouyinVideosMsg{
		ArrChannel: convertToRabbitMQArrChannel(msg.ArrChannel),
		IsScanFull: msg.IsScanFull,
		Group:      msg.Group,
	})
	if err != nil {
		uc.l.Errorf(ctx, "book.usecase.pubScanVideoOldTask.prod.pubRandomBookTask: %v", err)
		return err
	}

	return nil
}
