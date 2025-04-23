package middleware

import (
	"strings"

	"github.com/nguyenthanhtrung001/reup/pkg/jwt"
	"github.com/nguyenthanhtrung001/reup/pkg/response"

	"github.com/gin-gonic/gin"
)

func (m Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		if tokenString == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		payload, err := m.jwtManager.VerifyToken(tokenString)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// util.PrintJson(payload)

		// Check if empty payload.GroupID will return false
		if payload.GroupID == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		ctx = jwt.SetPayloadToContext(ctx, payload)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
