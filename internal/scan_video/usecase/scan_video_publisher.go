package usecase

import (
	"context"
	"reup/internal/scan_video/delivery/rabbitmq"
)

func (uc implUseCase) pubScanVideoOldTask(ctx context.Context, msg ScanDouyinVideosInput) error {
	uc.l.Info(ctx, "========== đẩy queue video old ==========")
	err := uc.prod.PubScanVideoOld(ctx, rabbitmq.ScanDouyinVideosMsg{
		Mid:        msg.Mid,
		SecUserID:  msg.SecUserID,
		VideoCount: msg.VideoCount,
		NewFlag:    msg.NewFlag,
		Group:      msg.Group,
		DomainAPI:  msg.DomainAPI,
	})
	if err != nil {
		uc.l.Errorf(ctx, "book.usecase.pubScanVideoOldTask.prod.pubRandomBookTask: %v", err)
		return err
	}

	return nil
}

func (uc implUseCase) pubScanVideoNewTask(ctx context.Context, msg ScanDouyinVideosInput) error {
	uc.l.Info(ctx, "========== đẩy queue video new ==========")
	err := uc.prod.PubScanVideoNew(ctx, rabbitmq.ScanDouyinVideosMsg{
		Mid:        msg.Mid,
		SecUserID:  msg.SecUserID,
		VideoCount: msg.VideoCount,
		NewFlag:    msg.NewFlag,
		Group:      msg.Group,
		DomainAPI:  msg.DomainAPI,
	})
	if err != nil {
		uc.l.Errorf(ctx, "book.usecase.pubScanVideoNewTask.prod.pubRandomBookTask: %v", err)
		return err
	}

	return nil
}
func (uc implUseCase) pubScanVideoManualTask(ctx context.Context, msg ScanDouyinVideosInput) error {
	uc.l.Info(ctx, "========== đẩy queue video manual ==========")
	err := uc.prod.PubScanVideoManual(ctx, rabbitmq.ScanDouyinVideosMsg{
		Mid:        msg.Mid,
		SecUserID:  msg.SecUserID,
		VideoCount: msg.VideoCount,
		NewFlag:    msg.NewFlag,
		Group:      msg.Group,
		DomainAPI:  msg.DomainAPI,
	})
	if err != nil {
		uc.l.Errorf(ctx, "book.usecase.pubScanVideoManualTask.prod.pubRandomBookTask: %v", err)
		return err
	}

	return nil
}
