package usecase

import (
	"context"

	"reup/internal/models"
	"reup/internal/scan_video/repository"
	"strings"
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

		// Cập nhật trạng thái cho từng BiliSpace
		for _, space := range douyinSpaces {
			err := uc.repo.UpdateBiliSpace(ctx, space.Mid, 0, 0)
			if err != nil {
				uc.l.Errorf(ctx, "Error updating douyin space for MID %s: %v", space.Mid, err)
			}
		}

		// Xử lý từng video Douyin
		for _, space := range douyinSpaces {
			secUserID := strings.Replace(space.DouyinLink, "https://www.douyin.com/user/", "", -1)

			// Thực hiện quét video Douyin
			data := ScanDouyinVideosInput{
				Mid:        space.Mid,
				SecUserID:  secUserID,
				VideoCount: 10,
				NewFlag:    false,
				Group:      1,
				DomainAPI:  cs.DomainAPI,
			}
			uc.pubScanVideoOldTask(ctx, data)

		}

		uc.l.Infof(ctx, "Finished processing computer: %s", cs.Computer)
	}

	uc.l.Infof(ctx, "========= SCAN DOUYIN VIDEOS COMPLETED =========")
}
