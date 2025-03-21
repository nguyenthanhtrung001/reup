package usecase

import (
	"context"
	"reup/internal/book/delivery/rabbitmq"
)

func (uc implUseCase) pubRandomBookTask(ctx context.Context, msg string) error {
	err := uc.prod.PubBookRandom(ctx, rabbitmq.BookMsg{
		Msg: "Trung đẹp trai vừa gửi tin ..." + msg,
	})
	if err != nil {
		uc.l.Errorf(ctx, "book.usecase.pubRandomBookTask.prod.pubRandomBookTask: %v", err)
		return err
	}

	return nil
}
