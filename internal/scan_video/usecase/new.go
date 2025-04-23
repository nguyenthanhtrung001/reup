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
	ProxyScan
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
	ScanDouyinVideoSheduler()
	ScanDouyinVideoFullPageSheduler()
	ScanDouyinVideoManualSheduler()
	HandleDouyinWebhook(c interface{}, data HandleDouyinWebhookInput) (int, error)
	SentScanDouyinVideos(spaceArr []ArrChannel, isFull bool) error
}

type ScanChannel interface {
	ScanChannel()
}

type Telegram interface {
	filterAndSendTelegram(ctx context.Context, videos []models.BiliVideo, froup int) error
}

type ProxyScan interface {
	DoneAllProxyScan(ctx context.Context) error
	DoneProxyScan(ctx context.Context, proxyIP string) error
	GetAllProxyScan(ctx context.Context) ([]models.ProxyScan, error)
	GetProxyScanRandom(ctx context.Context) (*models.ProxyScan, error)
	InsertProxyScan(ctx context.Context, ProxyIP string) error
	GetQuest(ctx context.Context, computer string) (*models.BiliVideo, map[string]interface{}, error)
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
