package usecase

import (
	"context"

	"reup/internal/models"
	"reup/internal/scan_video/repository"
	"strings"
)

func (uc implUseCase) ScanDouyinVideoFullPageSheduler() {
	ctx := context.Background()
	uc.l.Infof(ctx, "========= SCAN DOUYIN VIDEOS BY CHANNEL NEW... ==========")

	// Lấy danh sách các computer có tên khác rỗng và có đuôi _DY, _XIT
	computers, err := uc.repo.FindChannelsGroupedByComputer(ctx)
	if err != nil {
		uc.l.Errorf(ctx, "Error getting channels from DB: %v", err)
		return
	}
	settings, err := uc.repo.FindAllPriorityScanComputers(ctx)
	if err != nil {
		uc.l.Errorf(ctx, "Error getting settings from DB: %v", err)
		return
	}
	// Tìm cấu hình mặc định
	settingDefault := uc.findDefaultSetting(settings)
	if (settingDefault == models.PriorityScanComputer{}) {
		uc.l.Errorf(ctx, "No computer found with computer_name 'default'. Exiting...")
		return
	}

	// Tạo map lưu trữ số lượng scan cho mỗi computer
	scanNumbersMap := uc.createScanNumbersMap(settings)
	doaminMap := uc.createDomainMap(settings)
	// uc.l.Infof(ctx, "Created scan numbers map from settings.")

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
		uc.l.Infof(ctx, "Computer: %s, Scan Numbers: %d", computer.Computer, scanNumbers)
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
		uc.l.Infof(ctx, "Found %d users for computer %s.", len(users), cs.Computer)

		douyinLastScanIsZero := true
		// Thiết lập tham số tìm kiếm Douyin Spaces
		userScan := repository.FindDouyinOldSpacesOptions{
			ChannelUsernames:     users,
			ScanNumbers:          2,
			DouyinLastScanIsZero: &douyinLastScanIsZero,
		}

		// Lấy các BiliSpace từ repository
		douyinSpaces, err := uc.repo.FindDouyinOldSpaces(ctx, &userScan)
		if err != nil {
			uc.l.Errorf(ctx, "Error getting douyin spaces for computer %s: %v", cs.Computer, err)
			continue
		}
		uc.l.Infof(ctx, "Found %d douyin spaces for computer %s.", len(douyinSpaces), cs.Computer)

		// Cập nhật trạng thái cho từng BiliSpace
		for _, space := range douyinSpaces {
			err := uc.repo.UpdateBiliSpace(ctx, space.Mid, 0, 0)
			if err != nil {
				uc.l.Errorf(ctx, "Error updating douyin space for MID %s: %v", space.Mid, err)
			} else {
				uc.l.Infof(ctx, "Successfully updated douyin space for MID %s.", space.Mid)
			}
		}

		// Xử lý từng video Douyin
		for _, space := range douyinSpaces {
			secUserID := strings.Replace(space.DouyinLink, "https://www.douyin.com/user/", "", -1)
			uc.l.Infof(ctx, "Started video scan job for Douyin user %s.", secUserID)

			// Thực hiện quét video Douyin (có thể sử dụng job hoặc API để xử lý)
			// go scanDouyinVideosJob(space.Mid, secUserID, 10, false, 0, domainAPI)
			uc.l.Info(ctx, "Running job for user:", secUserID)

			data := ScanDouyinVideosInput{
				Mid:        space.Mid,
				SecUserID:  secUserID,
				VideoCount: 10,
				NewFlag:    true,
				Group:      0,
				DomainAPI:  cs.DomainAPI,
			}
			uc.pubScanVideoNewTask(ctx, data)
		}

		uc.l.Infof(ctx, "Finished processing computer: %s", cs.Computer)
	}

	uc.l.Infof(ctx, "========= SCAN DOUYIN VIDEOS COMPLETED =========")
}
