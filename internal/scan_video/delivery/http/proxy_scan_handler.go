package http

import (
	"fmt"
	"net/http"

	"github.com/nguyenthanhtrung001/reup/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h handler) DoneAllProxyScan(c *gin.Context) {
	ctx := c.Request.Context()
	err := h.uc.DoneAllProxyScan(ctx)
	if err != nil {
		h.l.Error(ctx, "error done all proxy scan", err)
		response.ErrorWithMap(c, err, nil)
		return
	}
	response.OK(c, "done all proxy scan")

}

// DoneProxyScan implements Handler.
func (h handler) DoneProxyScan(c *gin.Context) {
	ctx := c.Request.Context()
	var request struct {
		ProxyIP string `json:"proxy_ip"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.l.Error(ctx, "error binding proxy data", err)
		response.ErrorWithMap(c, err, nil)
		return
	}
	err := h.uc.DoneProxyScan(ctx, request.ProxyIP)
	if err != nil {
		h.l.Error(ctx, "error done proxy scan", err)
		response.ErrorWithMap(c, errDoneProxyScan, nil)
		return
	}
	response.OK(c, "done proxy scan")

}

// GetAllProxyScan implements Handler.
func (h handler) GetAllProxyScan(c *gin.Context) {
	ctx := c.Request.Context()
	proxyScans, err := h.uc.GetAllProxyScan(ctx)
	if err != nil {
		h.l.Error(ctx, "error get all proxy scan", err)
		response.ErrorWithMap(c, errGetAllProxyScan, nil)
		return
	}
	var proxyIPs []string
	for _, scan := range proxyScans {
		proxyIPs = append(proxyIPs, scan.ProxyIP)
	}

	response.OK(c, proxyIPs)

}

func (h handler) GetProxyScanRandom(c *gin.Context) {
	ctx := c.Request.Context()
	proxyScan, err := h.uc.GetProxyScanRandom(ctx)
	if err != nil {
		h.l.Error(ctx, "error get proxy scan random", err)
		response.ErrorWithMap(c, errGetProxyScanRandom, nil)
		return
	}
	response.OK(c, proxyScan.ProxyIP)
}

func (h handler) InsertProxyScan(c *gin.Context) {
	ctx := c.Request.Context()
	var request struct {
		ProxyIP string `json:"proxy_ip"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.l.Error(ctx, "error binding proxy data", err)
		response.ErrorWithMap(c, err, nil)
		return
	}
	err := h.uc.InsertProxyScan(ctx, request.ProxyIP)
	if err != nil {
		h.l.Error(ctx, "error insert proxy scan", err)
		response.ErrorWithMap(c, err, nil)
		return
	}
	response.OK(c, "insert proxy scan")
}

func (h handler) GetQuestHandler(c *gin.Context) {
	computer := c.Query("computer")
	if computer == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "computer param is required"})
		return
	}

	video, userProxyInfo, err := h.uc.GetQuest(c.Request.Context(), computer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Chuyển đổi dữ liệu video từ MongoDB sang định dạng mong muốn
	response := gin.H{
		"success": true,
		"data": gin.H{
			"video_id":           fmt.Sprintf("%v", video.VideoID), // Nếu VideoID là số thì chuyển thành string
			"mid":                video.Mid,
			"video_thumb":        video.VideoThumb,
			"upload_title":       video.UploadTitle,
			"upload_description": video.UploadDescription,
			"upload_keyword":     video.UploadKeyword,
			"duration":           video.Duration,
			"date_make":          video.DateMake.Format("2006-01-02 15:04:05"), // Đảm bảo đúng định dạng
			"count_get":          video.CountGet,
			"next":               video.Next,
			"download_fail":      video.DownloadFail,
			"type":               video.Type,
			"updated_at":         video.UpdatedAt.Format("2006-01-02T15:04:05.000000Z"), // Đảm bảo đúng định dạng
			"created_at":         video.CreatedAt.Format("2006-01-02T15:04:05.000000Z"),
			"download_url":       video.DownloadURL,
			"tags":               video.Tags,                                                    // Nếu không có tags, trả về null
			"video_link":         fmt.Sprintf("https://www.douyin.com/video/%v", video.VideoID), // Tạo link video Douyin
			"check_copyright":    0,                                                             // Bạn có thể điều chỉnh thêm các giá trị nếu cần
			"username":           userProxyInfo["username"],                                     // Bạn có thể thay bằng tên người dùng thực tế
			"proxy_ip":           userProxyInfo["proxy"],                                        // Cập nhật nếu cần
			"ads":                0,
		},
	}

	c.JSON(http.StatusOK, response)
}
