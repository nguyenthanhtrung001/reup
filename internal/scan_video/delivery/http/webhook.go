package http

import (
	"reup/internal/scan_video/usecase"

	"github.com/gin-gonic/gin"
)

func (h handler) HandleDouyinWebhook(c *gin.Context) {
	var data usecase.HandleDouyinWebhookInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	if data.ChannelID == "" || len(data.Posts) == 0 {
		c.JSON(400, gin.H{"error": "Missing required data"})
		return
	}

	// Call the use case to handle the business logic
	insertedCount, err := h.uc.HandleDouyinWebhook(c, data)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":        "Webhook received and processed successfully",
		"channel_id":     data.ChannelID,
		"inserted_count": insertedCount,
	})
}
