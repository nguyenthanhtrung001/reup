package http

import (
	userUseCase "github.com/nguyenthanhtrung001/reup/internal/user/usecase"
	"github.com/nguyenthanhtrung001/reup/pkg/log"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	UserHandlerHandler
}

type UserHandlerHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Profile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	UpdatePassword(c *gin.Context)

	// Admin User
	List(c *gin.Context)
	DetailAdmin(c *gin.Context)
	UpdateByAdmin(c *gin.Context)
	CreateApiKey(c *gin.Context)
}

type handler struct {
	l      log.Logger
	userUC userUseCase.UseCase
}

func New(l log.Logger, userUC userUseCase.UseCase) Handler {
	return handler{
		l:      l,
		userUC: userUC,
	}
}
