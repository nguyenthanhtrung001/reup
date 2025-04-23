package middleware

import (
	"github.com/nguyenthanhtrung001/reup/pkg/jwt"
	"github.com/nguyenthanhtrung001/reup/pkg/response"

	"github.com/gin-gonic/gin"
)

func (m Middleware) AuthKeyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		sc, ok := jwt.GetScopeFromContext(c.Request.Context())
		if !ok {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		if !sc.IsAdmin() && !sc.IsSuperAdmin() {
			response.PermissionDenied(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
