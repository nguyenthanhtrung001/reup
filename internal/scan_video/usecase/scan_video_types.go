package usecase

import (
	"github.com/nguyenthanhtrung001/reup/internal/models"
)

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

type ScanDouyinVideosInput struct {
	ArrChannel []ArrChannel
	IsScanFull bool
	Group      int
}
type ArrChannel struct {
	SpaceId   int64
	ChannelId []string
}

type ComputerSetting struct {
	Computer    string
	ScanNumbers int
	DomainAPI   string
}

type ScanResponse struct {
	Message     string `json:"message"`
	QueuedCount int    `json:"queued_count"`
}

// Struct để gửi request đến API
type ScanRequest struct {
	Spaces []ArrChannel `json:"channels"`
	IsFull bool         `json:"is_full"`
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

// convertToRabbitMQArrChannel converts a slice of map[string]interface{} to a slice of rabbitmq.ArrChannel
func convertToRabbitMQArrChannel2(channels []map[string]interface{}) []ArrChannel {
	result := make([]ArrChannel, len(channels))
	for i, ch := range channels {

		result[i] = ArrChannel{
			SpaceId:   ch["space_id"].(int64),
			ChannelId: []string{ch["channel_id"].(string)},
		}
	}
	return result
}
