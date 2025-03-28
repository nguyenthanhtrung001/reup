package usecase

import (
	"reup/internal/models"
	"time"
)

type DouyinResponse struct {
	AwemeList []VideoInfo `json:"aweme_list"`
	MaxCursor int         `json:"max_cursor"`
}

type DouyinVideo struct {
	VideoID       string `json:"video_id"`
	VideoThumb    string `json:"video_thumb"`
	UploadTitle   string `json:"upload_title"`
	UploadDesc    string `json:"upload_description"`
	UploadKeyword string `json:"upload_keyword"`
	Duration      int    `json:"duration"`
	DateMake      string `json:"date_make"`
	CountGet      int    `json:"count_get"`
	Next          string `json:"next"`
	DownloadFail  bool   `json:"download_fail"`
	Type          string `json:"type"`
	UpdatedAt     string `json:"updated_at"`
	CreatedAt     string `json:"created_at"`
}

type VideoInfo struct {
	AwemeID    string `json:"aweme_id"`
	Desc       string `json:"desc"`
	CreateTime int    `json:"create_time"`
	Duration   int    `binding:"duration"`
	Video      Video  `json:"video"`
}

type Video struct {
	Cover Cover `json:"cover"`
}
type Cover struct {
	URLList []string `json:"url_list"`
}

type ScanDouyinVideosInput struct {
	Mid        int64
	SecUserID  string
	VideoCount int
	NewFlag    bool
	Group      int
	DomainAPI  string
}

func intToTime(timestamp int64) time.Time {
	// Chuyển timestamp sang time.Time
	return time.Unix(timestamp, 0)
}

type ComputerSetting struct {
	Computer    string
	ScanNumbers int
	DomainAPI   string
}

type InputProducer struct {
	Mid       string
	SecUserId string
	Count     int
}

// Hàm tìm cấu hình mặc định
func (uc implUseCase) findDefaultSetting(settings []models.PriorityScanComputer) models.PriorityScanComputer {
	for _, setting := range settings {
		if setting.ComputerName == "default" {
			return setting
		}
	}
	return models.PriorityScanComputer{}
}

// Hàm tạo map scanNumbers từ cấu hình settings
func (uc implUseCase) createScanNumbersMap(settings []models.PriorityScanComputer) map[string]int {
	scanNumbersMap := make(map[string]int)
	for _, setting := range settings {
		if setting.ComputerName != "" {
			scanNumbersMap[setting.ComputerName] = setting.ScanNumbers
		}
	}
	return scanNumbersMap
}

// Hàm tạo map scanNumbers từ cấu hình settings
func (uc implUseCase) createDomainMap(settings []models.PriorityScanComputer) map[string]string {
	domainMap := make(map[string]string)
	for _, setting := range settings {
		if setting.ComputerName != "" {
			domainMap[setting.ComputerName] = setting.DomainAPI
		}
	}
	return domainMap
}
