package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/nguyenthanhtrung001/reup/internal/models"
	"github.com/nguyenthanhtrung001/reup/pkg/telegram"
)

func (uc implUseCase) filterAndSendTelegram(ctx context.Context, videos []models.BiliVideo, group int) error {
	if len(videos) == 0 {
		uc.l.Info(ctx, "No videos to process.")
		return nil
	}

	today, yesterday := getTodayAndYesterday()

	filteredVideos := make([]models.BiliVideo, 0, len(videos))
	for _, video := range videos { // fix fix
		if isVideoRecent(video.DateMake, today, yesterday) {
			filteredVideos = append(filteredVideos, video)
		}
	}

	if len(filteredVideos) == 0 {
		uc.l.Info(ctx, "No videos found for today or yesterday.")
		return nil
	}

	if len(filteredVideos) > 5 {
		filteredVideos = filteredVideos[:5]
	}

	err := uc.sendTelegram(ctx, filteredVideos, group)
	if err != nil {
		return fmt.Errorf("error sending telegram message: %v", err)
	}

	return nil
}

func isVideoRecent(makeDate time.Time, today time.Time, yesterday time.Time) bool {
	if makeDate.Year() == today.Year() && makeDate.YearDay() == today.YearDay() {
		return true
	}
	if makeDate.Year() == yesterday.Year() && makeDate.YearDay() == yesterday.YearDay() {
		return true
	}
	return false
}

func getTodayAndYesterday() (today time.Time, yesterday time.Time) {
	today = time.Now()
	yesterday = today.AddDate(0, 0, -1)
	return today, yesterday
}

func (uc implUseCase) sendTelegram(ctx context.Context, videos []models.BiliVideo, group int) error {
	// Tạo mảng để chứa thông tin video
	resArray := make(map[string]interface{})
	count := 1
	chatId := 0

	switch group {
	case 1:
		chatId = int(uc.teleChat.GroupChat1)
	case 2:
		chatId = int(uc.teleChat.GroupChat2)
	case 3:
		chatId = int(uc.teleChat.GroupChat3)
	}
	uc.l.Info(ctx, "chat id:", chatId)

	// Lặp qua các video và xây dựng thông điệp
	for _, video := range videos {
		// Tạo thông tin video với định dạng ngày "dd/mm/yyyy HH:mm:ss"
		notiVideoInfo := map[string]interface{}{
			"date":     video.DateMake.Format("02/01/2006 15:04:05"), // Định dạng ngày thành "dd/mm/yyyy HH:mm:ss"
			"video_id": fmt.Sprintf("https://www.douyin.com/video/%s", video.VideoID),
		}

		// Thêm vào mảng kết quả với key "Video-count"
		resArray[fmt.Sprintf("Video-%d", count)] = notiVideoInfo
		count++
	}

	// Tạo thông điệp text từ resArray
	message := fmt.Sprintf("Video Information:\n%s", formatResArrayToString(resArray))
	// Gửi thông điệp qua Telegram
	_, err := uc.tele.SendMessage(int64(chatId), message, telegram.HTMLMode)
	if err != nil {
		return fmt.Errorf("error sending message to Telegram: %v", err)
	}

	return nil
}

func formatResArrayToString(resArray map[string]interface{}) string {
	var result string

	if resArray == nil || len(resArray) == 0 {
		return "No video information available"
	}

	for key, videoInfo := range resArray {
		info, ok := videoInfo.(map[string]interface{})
		if !ok {
			result += fmt.Sprintf("\nInvalid data for key: %s\n", key)
			continue
		}

		videoID, okVideoID := info["video_id"].(string)
		date, okDate := info["date"].(string)

		result += fmt.Sprintf("\n%s:\n", key)

		if okVideoID {
			result += fmt.Sprintf("  Video ID: %s\n", videoID)
		} else {
			result += "  Video ID: Not available\n"
		}

		// Đảm bảo rằng date được phân tích đúng theo định dạng "dd/mm/yyyy HH:mm:ss"
		if okDate {
			parsedDate, err := time.Parse("02/01/2006 15:04:05", date) // Định dạng "dd/mm/yyyy HH:mm:ss"
			if err == nil {
				result += fmt.Sprintf("  Date: %s\n", parsedDate.Format("02/01/2006 15:04:05"))
			} else {
				result += "  Date: Invalid format\n"
			}
		} else {
			result += "  Date: Not available\n"
		}
	}

	return result
}
