package middleware

import (
	"github.com/nguyenthanhtrung001/reup/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	system_maintenance = false
)

func (m Middleware) SystemMaintenance() gin.HandlerFunc {
	return func(c *gin.Context) {
		if system_maintenance {
			response.SystemUnderMaintenance(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
