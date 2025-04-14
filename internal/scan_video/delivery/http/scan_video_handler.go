package http

import (
	"context"
	"net/http"
	"reup/internal/scan_video/usecase"

	"github.com/gin-gonic/gin"
)

func (h handler) scanDouyinVideosHandler(c *gin.Context) {
	// Tạo đối tượng input cho ScanDouyinVideos
	input := usecase.ScanDouyinVideosInput{
		Mid:        12345678,                                                  // Ví dụ mid
		SecUserID:  "MS4wLjABAAAA2I6OEC1NbqDri0clGwXp7o46Fx5ivPXtA8-brwBURhM", // Ví dụ sec_user_id
		VideoCount: 10,                                                        // Ví dụ số lượng video cần quét
		NewFlag:    true,                                                      // Quét video mới
		Group:      99,                                                        // Ví dụ nhóm
		DomainAPI:  "103.42.56.42:3000",                                       // Ví dụ domain API
	}
	abc, _ := h.uc.ScanDouyinVideos(context.Background(), input)

	// // videos, err := h.uc.ScanDouyinVideos(context.Background(), input)
	// // if err != nil {
	// // 	// Nếu có lỗi, trả về mã lỗi HTTP và thông báo lỗi
	// // 	c.JSON(http.StatusInternalServerError, gin.H{
	// // 		"error": fmt.Sprintf("Error scanning Douyin videos: %v", err),
	// // 	})
	// // 	return
	// // }

	// h.uc.ScanDouyinVideoSheduler()

	// Chuyển kết quả videos thành JSON và gửi về cho client
	c.JSON(http.StatusOK, abc)
}
