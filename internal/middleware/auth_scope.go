package middleware

import (
	"book-store/pkg/jwt"
	"book-store/pkg/response"

	"github.com/gin-gonic/gin"
)

func (m Middleware) AuthScope() gin.HandlerFunc {
	return func(c *gin.Context) {
		scopeString := c.GetHeader("Scope")
		if scopeString == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		scope, err := jwt.ParseScopeHeader(scopeString)
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
