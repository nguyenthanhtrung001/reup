package job

import (
	"context"
)

func (h Handler) testJobBook() {
	ctx := context.Background()
	h.l.Info(ctx, "test job")
	book, err := h.uc.GetRandomBookTask()
	if err != nil {
		h.l.Errorf(ctx, "job.handler.test: %v", err)
	}

	h.l.Info(ctx, book.Title)

}
