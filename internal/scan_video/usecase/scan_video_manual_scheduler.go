package usecase

import (
	"context"
	"strings"

	"github.com/nguyenthanhtrung001/reup/internal/models"
	"github.com/nguyenthanhtrung001/reup/internal/scan_video/repository"
)

func (uc implUseCase) ScanDouyinVideoManualSheduler() {
	ctx := context.Background()
	uc.l.Info(ctx, "======================== Start manual scan ===========================")

	// Lấy cấu hình từ DB
	setting, err := uc.repo.GetFirstRecord(ctx)
	if err != nil {
		uc.l.Errorf(ctx, "Error getting settings from DB: %v", err)
		return
	}

	var douyinSpaces []models.BiliSpace

	// Nếu trạng thái quét video thủ công là 1 (bật)
	if setting.DouyinManualScanStatus == 1 {
		uc.l.Info(ctx, "Manual scan is enabled, fetching Douyin spaces...")

		// Lấy danh sách Douyin spaces từ repository
		users := []string{"manual_group2@gmail.com"}
		userScan := repository.FindDouyinOldSpacesOptions{
			ChannelUsernames: users,
			ScanNumbers:      200, // Giới hạn số lượng quét
		}

		// Lấy các Douyin spaces từ repository
		douyinSpaces, err = uc.repo.FindDouyinOldSpaces(ctx, &userScan)
		if err != nil {
			uc.l.Errorf(ctx, "Error getting Douyin spaces: %v", err)
			return
		}

		uc.l.Infof(ctx, "Fetched %d Douyin spaces successfully.", len(douyinSpaces))

	} else {
		uc.l.Info(ctx, "========= SCAN (MANUALLY) DOUYIN VIDEOS BY CHANNEL IS DISABLED... ==========")
		douyinSpaces = []models.BiliSpace{} // Nếu không quét, gán slice rỗng
	}

	// Tạo các mảng để phân loại quét tất cả video và video mới
	var scanAllVideosArr []map[string]interface{}
	var scanNewVideosArr []map[string]interface{}

	// Xử lý từng Douyin space và phân loại
	for _, space := range douyinSpaces {
		// Trích xuất channel_id từ DouyinLink
		channelID := strings.Replace(space.DouyinLink, "https://www.douyin.com/user/", "", -1)
		spaceArr := map[string]interface{}{
			"space_id":   space.Mid,
			"channel_id": channelID,
		}

		// Lọc ra 2 nhóm: quét tất cả video và quét video mới
		if space.DouyinLastScan == nil || *space.DouyinLastScan == 0 || space.CountVideo == 0 {
			scanAllVideosArr = append(scanAllVideosArr, spaceArr)
		} else {
			scanNewVideosArr = append(scanNewVideosArr, spaceArr)
		}
	}

	// Log mảng gửi đi
	if len(scanAllVideosArr) > 0 {
		uc.l.Infof(ctx, "Dispatching job to scan all videos with %d entries.", len(scanAllVideosArr))
		// Gửi tác vụ quét tất cả video
		input := ScanDouyinVideosInput{
			ArrChannel: convertToArrayDataSendScan(scanAllVideosArr),
			IsScanFull: true, // Cập nhật flag quét đầy đủ tùy theo yêu cầu
			Group:      1,
		}

		err = uc.SentScanDouyinVideos(input.ArrChannel, input.IsScanFull)

		if err != nil {
			uc.l.Errorf(ctx, "Error publishing to RabbitMQ for scan all videos: %v", err)
		}
	}

	if len(scanNewVideosArr) > 0 {
		uc.l.Infof(ctx, "Dispatching job to scan new videos with %d entries.", len(scanNewVideosArr))
		// Gửi tác vụ quét video mới
		input := ScanDouyinVideosInput{
			ArrChannel: convertToArrayDataSendScan(scanNewVideosArr),
			IsScanFull: false, // Cập nhật flag quét đầy đủ tùy theo yêu cầu
			Group:      1,
		}

		err = uc.SentScanDouyinVideos(input.ArrChannel, input.IsScanFull)

		if err != nil {
			uc.l.Errorf(ctx, "Error publishing to RabbitMQ for scan new videos: %v", err)
		}
	}

	uc.l.Info(ctx, "======================== End of manual scan ===========================")
}
