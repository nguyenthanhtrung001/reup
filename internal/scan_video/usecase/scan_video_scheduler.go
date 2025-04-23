package usecase

import (
	"context"

	"strings"

	"github.com/nguyenthanhtrung001/reup/internal/models"
	"github.com/nguyenthanhtrung001/reup/internal/scan_video/repository"
)

func (uc implUseCase) ScanDouyinVideoSheduler() {
	ctx := context.Background()
	uc.l.Infof(ctx, "========= SCAN DOUYIN VIDEOS BY CHANNEL OLD... ==========")

	// Lấy danh sách các computer
	computers, err := uc.repo.FindChannelsGroupedByComputer(ctx)
	if err != nil {
		uc.l.Errorf(ctx, "Error getting channels from DB: %v", err)
		return
	}

	// Lấy tất cả các cấu hình từ PriorityScanComputers
	settings, err := uc.repo.FindAllPriorityScanComputers(ctx)
	if err != nil {
		uc.l.Errorf(ctx, "Error getting settings from DB: %v", err)
		return
	}

	// Tìm cấu hình mặc định
	settingDefault := uc.findDefaultSetting(settings)
	if (settingDefault == models.PriorityScanComputer{}) {
		uc.l.Errorf(ctx, "No computer found with computer_name 'default'. Exiting...", computers)
		return
	}

	// Tạo map lưu trữ số lượng scan cho mỗi computer
	scanNumbersMap := uc.createScanNumbersMap(settings)
	doaminMap := uc.createDomainMap(settings)

	// Xử lý thông tin cấu hình của các computer
	computerSettings := make(map[string]ComputerSetting)
	for _, computer := range computers {
		scanNumbers := scanNumbersMap[computer.Computer]
		domain := doaminMap[computer.Computer]
		computerSettings[computer.Computer] = ComputerSetting{
			Computer:    computer.Computer,
			ScanNumbers: scanNumbers,
			DomainAPI:   domain,
		}
	}

	// Xử lý từng computer và quét các Douyin spaces
	for _, cs := range computerSettings {
		uc.l.Infof(ctx, "Processing computer: %s with %d scan numbers.", cs.Computer, cs.ScanNumbers)

		// Lấy danh sách usernames từ repository
		users, err := uc.repo.FindUsernamesByComputer(ctx, cs.Computer)
		if err != nil {
			uc.l.Errorf(ctx, "Error getting users for computer %s: %v", cs.Computer, err)
			continue
		}

		douyinLastScanIsZero := false
		// Thiết lập tham số tìm kiếm Douyin Spaces
		userScan := repository.FindDouyinOldSpacesOptions{
			ChannelUsernames:     users,
			ScanNumbers:          int64(cs.ScanNumbers),
			DouyinLastScanIsZero: &douyinLastScanIsZero,
		}

		// Lấy các BiliSpace từ repository
		douyinSpaces, err := uc.repo.FindDouyinOldSpaces(ctx, &userScan)
		if err != nil {
			uc.l.Errorf(ctx, "Error getting douyin spaces for computer %s: %v", cs.Computer, err)
			continue
		}

		// Cập nhật trạng thái cho từng BiliSpace gồm :douyin_last_scan, last_scan
		for _, space := range douyinSpaces {
			err := uc.repo.UpdateBiliSpace(ctx, space.Mid, 0, 0)
			if err != nil {
				uc.l.Errorf(ctx, "Error updating douyin space for MID %s: %v", space.Mid, err)
			}
		}

		// Xử lý từng video Douyin
		spaceArr := make([]map[string]interface{}, len(douyinSpaces))
		for i, space := range douyinSpaces {
			// Trích xuất channel_id từ DouyinLink
			channelID := strings.Replace(space.DouyinLink, "https://www.douyin.com/user/", "", -1)
			spaceArr[i] = map[string]interface{}{
				"space_id":   space.Mid,
				"channel_id": channelID,
			}
		}
		err = uc.pubScanVideoOldTask(ctx, ScanDouyinVideosInput{
			ArrChannel: convertToRabbitMQArrChannel2(spaceArr),
			IsScanFull: false, // Cập nhật flag quét đầy đủ tùy theo yêu cầu
			Group:      1,     // Có thể thay đổi nhóm nếu cần
		})
		if err != nil {
			uc.l.Errorf(ctx, "Error publishing to RabbitMQ for computer %s: %v", err)
			continue
		}

		uc.l.Infof(ctx, "Finished processing computer: %s", cs.Computer)
	}

	uc.l.Infof(ctx, "========= SCAN DOUYIN VIDEOS COMPLETED =========")
}
