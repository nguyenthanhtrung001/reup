package usecase

import (
	"context"
	"reup/internal/models"
	"reup/internal/scan_video/repository"
	"strings"
)

func (uc implUseCase) ScanDouyinVideoManualSheduler() {

	ctx := context.Background()
	uc.l.Info(ctx, "======================== Start manual scan ===========================")

	setting, err := uc.repo.GetFirstRecord(ctx)
	if err != nil {
		uc.l.Errorf(ctx, "Error getting settings from DB: %v", err)
		return
	}

	var douyinSpaces []models.BiliSpace

	// Nếu trạng thái quét video thủ công là 1 (bật)
	if setting.DouyinManualScanStatus == 1 {
		uc.l.Info(ctx, "Manual scan is enabled, fetching Douyin spaces...")
		users := []string{"manual_group2@gmail.com"}
		userScan := repository.FindDouyinOldSpacesOptions{
			ChannelUsernames: users,
			ScanNumbers:      15,
		}

		// Lấy danh sách các Douyin spaces
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

	// Cập nhật trạng thái cho từng BiliSpace
	for _, space := range douyinSpaces {
		uc.l.Infof(ctx, "Updating Douyin space for MID: %s", space.Mid)
		err := uc.repo.UpdateBiliSpace(ctx, space.Mid, 0, 0)
		if err != nil {
			uc.l.Errorf(ctx, "Error updating Douyin space for MID %s: %v", space.Mid, err)
		} else {
			uc.l.Infof(ctx, "Successfully updated Douyin space for MID %s.", space.Mid)
		}
	}

	// Xử lý quét video Douyin
	for _, space := range douyinSpaces {
		secUserID := strings.Replace(space.DouyinLink, "https://www.douyin.com/user/", "", -1)

		uc.l.Infof(ctx, "Preparing scan for user: %s (SecUserID: %s)", space.DouyinLink, secUserID)

		// Nếu DouyinLastScan là nil hoặc không có video
		videoCount := 100
		newFlag := true
		if space.DouyinLastScan != nil && *space.DouyinLastScan != 0 && space.CountVideo > 0 {
			videoCount = 10
			newFlag = false
			uc.l.Infof(ctx, "Douyin last scan found, changing video count to %d", videoCount)
		}
		if newFlag == true {
			uc.l.Infof(ctx, "Douyin last scan found.........scan new")
		} else {
			uc.l.Infof(ctx, "Douyin last scan found.........scan old")
		}
		// Tạo đầu vào và gửi tác vụ quét video
		data := ScanDouyinVideosInput{
			Mid:        space.Mid,
			SecUserID:  secUserID,
			VideoCount: videoCount,
			NewFlag:    newFlag,
			Group:      0,
			DomainAPI:  "103.42.56.42:3000",
		}
		uc.l.Infof(ctx, "Sending video scan task for SecUserID: %s, VideoCount: %d", secUserID, videoCount)
		uc.pubScanVideoManualTask(ctx, data)

		uc.l.Info(ctx, "Job started for user:", secUserID)
	}

	uc.l.Info(ctx, "======================== End of manual scan ===========================")
}
