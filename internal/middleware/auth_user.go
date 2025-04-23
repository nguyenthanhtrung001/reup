package middleware

import (
	"github.com/nguyenthanhtrung001/reup/pkg/jwt"
	"github.com/nguyenthanhtrung001/reup/pkg/response"

	"github.com/gin-gonic/gin"
)

func (m Middleware) AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy payload từ context
		payload, ok := jwt.GetPayloadFromContext(c.Request.Context())

		if !ok {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// Khởi tạo scope từ payload
		sc := jwt.NewScope(payload)

		// Lấy thông tin người dùng
		user, err := m.userUC.GetUserById(c, sc, sc.UserID)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// Kiểm tra quyền người dùng
		if !user.IsUser() {
			response.PermissionDenied(c)
			c.Abort()
			return
		}

		// Tiếp tục xử lý request
		c.Next()
	}
}
