package schedule

import (
	"context"
	"os"
	"os/signal"

	"syscall"

	scanVideoJob "github.com/nguyenthanhtrung001/reup/internal/scan_video/delivery/job"
	scanVideoProd "github.com/nguyenthanhtrung001/reup/internal/scan_video/delivery/rabbitmq/producer"
	scanVideoMongo "github.com/nguyenthanhtrung001/reup/internal/scan_video/repository/mongo"
	scanVideoUsecase "github.com/nguyenthanhtrung001/reup/internal/scan_video/usecase"
	"github.com/nguyenthanhtrung001/reup/pkg/cron"
	"github.com/nguyenthanhtrung001/reup/pkg/jwt"
	"github.com/nguyenthanhtrung001/reup/pkg/telegram"
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

	chatIDs := telegram.ChatIDs{
		ReportBug:     s.telegram.ChatIDs.ReportBug,
		ReportPayment: s.telegram.ChatIDs.ReportPayment,
	}
	telegram := telegram.New(s.telegram.BotKey, chatIDs)

	// Producers
	scanVideoProd := scanVideoProd.New(s.l, s.conn)
	if err := scanVideoProd.Run(); err != nil {
		return err
	}

	// khai báo usecase sử dụng trong job
	scanVideoMongo := scanVideoMongo.New(s.l, s.db, jwt.JWTMaker{})
	scanVideoUsecase := scanVideoUsecase.New(s.l, scanVideoMongo, s.encrypter, scanVideoProd, telegram, scanVideoUsecase.TeleChat{
		NotifiChatID: s.telegram.ReportPayment,
	})

	jobHandler := []interface {
		Register() []cron.JobInfo
	}{
		// đăng ký job tại đây
		scanVideoJob.New(s.l, scanVideoUsecase, s.cron),
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
