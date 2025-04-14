package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"reup/internal/models"
	"reup/pkg/curl"
	"time"
)

func (uc implUseCase) ScanDouyinVideos(ctx context.Context, input ScanDouyinVideosInput) ([]models.BiliVideo, error) {
	// Lấy URL tùy thuộc vào Domain API
	url := fmt.Sprintf("http://%s/v1/social/douyin/web/aweme/post", input.DomainAPI)

	// Khởi tạo biến cần thiết
	var videos []models.BiliVideo
	var maxCursor = 0
	var count = 0

	// Kiểm tra flag để quyết định quét video mới hay cũ
	if input.NewFlag {
		// Quét video mới
		for count < 5 {
			// Gửi yêu cầu API
			resp, err := uc.sendRequest(url, input.SecUserID, input.VideoCount, maxCursor)
			if err != nil {
				return nil, err
			}
			if resp == nil || resp.MaxCursor == 0 {
				return nil, nil
			}

			// Xử lý dữ liệu video từ API
			videoItems := uc.processVideos(resp.AwemeList, input.Mid)

			// Thêm video vào danh sách
			videos = append(videos, videoItems...)

			// Cập nhật maxCursor cho lần request sau
			if resp.MaxCursor != 0 { // Kiểm tra nếu MaxCursor không phải là 0
				maxCursor = resp.MaxCursor
			}

			count++
		}
	} else {
		// Quét video cũ (chỉ một lần request)
		resp, err := uc.sendRequest(url, input.SecUserID, input.VideoCount, 0)
		if err != nil {
			return nil, err
		}

		// Xử lý dữ liệu video từ API
		videos = uc.processVideos(resp.AwemeList, input.Mid)
	}

	// Chèn video vào cơ sở dữ liệu
	if err := uc.repo.InsertBiliVideos(ctx, videos); err != nil {
		return nil, fmt.Errorf("error inserting videos: %v", err)
	}
	uc.l.Infof(ctx, "Inserted %d new videos into the database.", len(videos))

	// Cập nhật BiliSpace
	if err := uc.repo.UpdateBiliSpace(ctx, input.Mid, 1, 0); err != nil {
		return nil, fmt.Errorf("error updating bili space: %v", err)
	}
	if input.Group != 0 {
		err := uc.filterAndSendTelegram(ctx, videos, input.Group)
		if err == nil {
			uc.l.Infof(ctx, "Gửi tn thanh conmg", input.Mid)
		} else {
			uc.l.Infof(ctx, "Gửi tn thất bại", err)
		}
	}

	uc.l.Infof(ctx, "BiliSpace for MID %d updated successfully.", input.Mid)
	return videos, nil
}

// sendRequest gửi HTTP request tới Douyin API để lấy dữ liệu video
func (uc implUseCase) sendRequest(url, secUserID string, videoCount int, maxCursor int) (*DouyinResponse, error) {
	// Tạo body cho request
	reqBody := map[string]interface{}{
		"sec_user_id": secUserID,
		"count":       videoCount,
		"max_cursor":  maxCursor,
	}

	// Gửi yêu cầu POST sử dụng hàm Post từ package curl
	respStr, err := curl.Post(url, nil, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Douyin API: %v", err)
	}

	// Giải mã JSON response
	var response DouyinResponse
	err = json.Unmarshal([]byte(respStr), &response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response body: %v. Response: %s", err, respStr)
	}

	return &response, nil
}

// processVideos xử lý danh sách video trả về và chuẩn bị dữ liệu để lưu
func (uc implUseCase) processVideos(awemeList []VideoInfo, mid int64) []models.BiliVideo {
	var videos []models.BiliVideo
	for _, video := range awemeList {
		// Nếu video đã tồn tại, bỏ qua
		if uc.checkIfVideoExists(video.AwemeID) {
			continue
		}

		// Chuẩn bị dữ liệu video
		videoData := models.BiliVideo{
			VideoID:           video.AwemeID,
			Mid:               mid,
			VideoThumb:        video.Video.Cover.URLList[0],
			UploadTitle:       video.Desc,
			Duration:          video.Duration / 1000,
			UploadDescription: video.Desc,
			Type:              "douyin",
			UpdatedAt:         time.Now(),
			CreatedAt:         time.Now(),
			DateMake:          intToTime(int64(video.CreateTime)),
		}
		videos = append(videos, videoData)
	}
	return videos
}

// checkIfVideoExists kiểm tra video có tồn tại trong cơ sở dữ liệu không
func (uc implUseCase) checkIfVideoExists(videoID string) bool {
	return uc.repo.CheckDouyinVideoExists(context.Background(), videoID)
}
