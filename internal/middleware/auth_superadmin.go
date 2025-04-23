package middleware

import (
	"github.com/nguyenthanhtrung001/reup/pkg/jwt"
	"github.com/nguyenthanhtrung001/reup/pkg/response"

	"github.com/gin-gonic/gin"
)

func (m Middleware) AuthSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, ok := jwt.GetPayloadFromContext(c.Request.Context())
		if !ok {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		sc := jwt.NewScope(payload)
		user, err := m.userUC.GetUserById(c, sc, sc.UserID)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		if !user.IsSuperAdmin() {
			response.PermissionDenied(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
