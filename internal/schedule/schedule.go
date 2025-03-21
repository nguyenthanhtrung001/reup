package schedule

import (
	"context"
	"os"
	"os/signal"
	bookJob "reup/internal/book/delivery/job"
	bookProd "reup/internal/book/delivery/rabbitmq/producer"
	bookMongo "reup/internal/book/repository/mongo"
	bookUsecase "reup/internal/book/usecase"
	"reup/pkg/cron"
	"reup/pkg/jwt"
	"syscall"
)

func (s Scheduler) Start() error {
	ctx := context.Background()
	s.l.Info(ctx, "Starting scheduler")
	if err := s.registerJobs(); err != nil {
		return err
	}
	go func() {
		s.l.Info(ctx, "Starting cron job")
		s.cron.Start()
	}()
	// tạo 1 channel nhận tối đa 1 tín hiệu từ HDH (Signal)
	quit := make(chan os.Signal, 1)
	// đăng ký tín hiệu để dừng chương trình
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	s.l.Info(ctx, "Stopping scheduler")
	s.cron.Stop()
	return nil

}
func (s Scheduler) registerJobs() error {
	// tạo wrapper cho các job
	s.cron.SetFuncWrapper(s.jobWrapper)

	// Producers
	bookProd := bookProd.New(s.l, s.conn)
	if err := bookProd.Run(); err != nil {
		return err
	}

	// khai báo usecase sử dụng trong job
	bookMongo := bookMongo.New(s.l, s.db, jwt.JWTMaker{})
	bookUsecase := bookUsecase.New(s.l, bookMongo, s.encrypter, bookProd)

	s.l.Info(context.Background(), bookUsecase)

	jobHandler := []interface {
		Register() []cron.JobInfo
	}{
		// đăng ký job tại đây
		bookJob.New(s.l, bookUsecase, s.cron),
	}
	for _, job := range jobHandler {
		infos := job.Register()
		for _, info := range infos {
			if err := s.cron.AddJob(info); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s Scheduler) jobWrapper(f cron.HandleFunc) {
	// đống gói toàn bộ job để gôm lỗi
	defer func() {
		if err := recover(); err != nil {
			s.l.Error(context.Background(), "Recored form panic: %v", err)
		}
	}()
	f()
}
