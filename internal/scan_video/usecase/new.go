package usecase

import (
	"context"
	prod "reup/internal/scan_video/delivery/rabbitmq/producer"
	"reup/internal/scan_video/repository"

	"reup/internal/models"
	"reup/pkg/encrypter"
	"reup/pkg/log"
	"reup/pkg/telegram"
)

type UseCase interface {
	DouyinSpaceUC
	PriorityScanComputerUC
	ChannelUC
	DouyinVideoUC
	ScanVideo
	Telegram
}

type DouyinSpaceUC interface {
	FindDouyinOldSpaces(ctx context.Context, input *FindDouyinOldSpacesInput) ([]models.BiliSpace, error)
	CreateBiliSpace(ctx context.Context, input CreateBiliSpaceInput) (*models.BiliSpace, error)
	UpdateBiliSpace(ctx context.Context, mid int64) error
}

type PriorityScanComputerUC interface {
	FindAllPriorityScanComputers(ctx context.Context) ([]models.PriorityScanComputer, error)
}

type ChannelUC interface {
	GetChannelGroupsByComputer(ctx context.Context) ([]repository.ChannelGroup, error)
	CheckChannelExists(ctx context.Context, channelID int64) (bool, error)
}
type DouyinVideoUC interface {
	CheckDouyinVideoExists(ctx context.Context, videoID string) bool
	CreateDouyinVideo(ctx context.Context, input CreateBiliVideoInput) (*models.BiliVideo, error)
}
type ScanVideo interface {
	ScanDouyinVideos(ctx context.Context, input ScanDouyinVideosInput) ([]models.BiliVideo, error)
	ScanDouyinVideoSheduler()
	ScanDouyinVideoFullPageSheduler()
	ScanDouyinVideoManualSheduler()
}

type ScanChannel interface {
	ScanChannel()
}

type Telegram interface {
	filterAndSendTelegram(ctx context.Context, videos []models.BiliVideo, froup int) error
}
type implUseCase struct {
	l         log.Logger
	repo      repository.Repository
	encrypter encrypter.Encrypter
	prod      prod.Producer
	tele      telegram.Telegram
	teleChat  TeleChat
}
type TeleChat struct {
	NotifiChatID int64
	GroupChat1   int64
	GroupChat2   int64
	GroupChat3   int64
}

func New(l log.Logger,
	repo repository.Repository,
	encrypter encrypter.Encrypter,
	prod prod.Producer,
	tele telegram.Telegram,
	telChat TeleChat) UseCase {
	return implUseCase{
		l:         l,
		repo:      repo,
		encrypter: encrypter,
		prod:      prod,
		tele:      tele,
		teleChat:  telChat,
	}
}
