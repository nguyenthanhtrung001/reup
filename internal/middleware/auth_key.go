package middleware

import (
	"github.com/nguyenthanhtrung001/reup/pkg/jwt"
	"github.com/nguyenthanhtrung001/reup/pkg/response"

	"github.com/gin-gonic/gin"
)

type RestApiRequest struct {
	Key string `json:"key" form:"key" binding:"required"`
}

func (m Middleware) AuthKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		keyString := c.GetHeader("Key")
		if keyString == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		scope, err := jwt.ParseApiKey(keyString, m.encrypter)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		ctx := c.Request.Context()
		ctx = jwt.SetScopeToContext(ctx, scope)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
