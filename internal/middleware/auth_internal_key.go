package middleware

import (
	"github.com/nguyenthanhtrung001/reup/pkg/response"

	"github.com/gin-gonic/gin"
)

func (m Middleware) AuthInternalKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		secKeyEncoded := c.GetHeader("Internal-Key")
		if secKeyEncoded == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		secKey, err := m.encrypter.Decrypt(secKeyEncoded)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		if secKey != m.secretKey {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
