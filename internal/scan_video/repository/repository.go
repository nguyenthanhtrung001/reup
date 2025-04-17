package repository

import (
	"context"
	"reup/internal/models"
)

type Repository interface {
	BiliSpaceRepo
	PriorityScanComputerRepo
	ChannelRepo
	BiliVideoRepo
	AppSettingRepo
	ScanChannelRepo
	VideoWebHookRepo
	ProxyScanRepo
}

type BiliSpaceRepo interface {
	FindDouyinOldSpaces(ctx context.Context, opts *FindDouyinOldSpacesOptions) ([]models.BiliSpace, error)
	CreateBiliSpace(ctx context.Context, input CreateBiliSpaceInput) (*models.BiliSpace, error)
	UpdateBiliSpaceF(ctx context.Context, mid int64, updatedFields map[string]interface{}) error
	UpdateBiliSpace(ctx context.Context, mid int64, douyinWaitScan int, scanError int) error
}

type PriorityScanComputerRepo interface {
	FindAllPriorityScanComputers(ctx context.Context) ([]models.PriorityScanComputer, error)
}

type ChannelRepo interface {
	FindChannelsGroupedByComputer(ctx context.Context) ([]ChannelGroup, error)
	CheckChannelExists(ctx context.Context, channelID int64) (bool, error)
	FindUsernamesByComputer(ctx context.Context, computerName string) ([]string, error)
}

type BiliVideoRepo interface {
	CheckDouyinVideoExists(ctx context.Context, videoID string) bool
	CreateDouyinVideo(ctx context.Context, input CreateBiliVideoInput) (*models.BiliVideo, error)
	InsertBiliVideos(ctx context.Context, videos []models.BiliVideo) error
}

type AppSettingRepo interface {
	GetFirstRecord(ctx context.Context) (*models.AppSetting, error)
}

type ScanChannelRepo interface {
	FindDouyinScanChannelBySecUID(ctx context.Context, secUID string) (*models.DouyinScanChannel, error)
	InsertDouyinScanChannels(ctx context.Context, channels []models.DouyinScanChannel) error
}

type VideoWebHookRepo interface {
	GetExistingVideoIds(ctx context.Context, videoIds []string) (map[string]struct{}, error)
	GetBiliSpaceByDouyinLink(ctx context.Context, douyinLink string) (*models.BiliSpace, error)
}
type ProxyScanRepo interface {
	GetProxyScanRandom(ctx context.Context) (*models.ProxyScan, error)
	GetAllProxyScan(ctx context.Context) ([]models.ProxyScan, error)
	DoneProxyScan(ctx context.Context, proxyIP string) error
	DoneAllProxyScan(ctx context.Context) error
	InsertProxyScan(ctx context.Context, proxyIP string) error
}
