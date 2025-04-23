package middleware

import (
	"log"

	"github.com/nguyenthanhtrung001/reup/pkg/response"

	"github.com/nguyenthanhtrung001/reup/pkg/telegram"

	"github.com/gin-gonic/gin"
)

func Recovery(t telegram.Telegram, chatBugID int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[Panic Recovered] %v\n", err)
				log.Printf("Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
				response.PanicError(c, err, t, chatBugID)
			}
		}()
		c.Next()
	}
}
