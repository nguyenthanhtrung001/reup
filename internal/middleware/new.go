package middleware

import (
	userUseCase "github.com/nguyenthanhtrung001/reup/internal/user/usecase"
	pkgCrt "github.com/nguyenthanhtrung001/reup/pkg/encrypter"
	"github.com/nguyenthanhtrung001/reup/pkg/jwt"
	"github.com/nguyenthanhtrung001/reup/pkg/log"
)

type Middleware struct {
	l          log.Logger
	jwtManager jwt.Maker
	userUC     userUseCase.UseCase
	encrypter  pkgCrt.Encrypter
	secretKey  string
}

func New(l log.Logger, jwtManager jwt.Maker, userUC userUseCase.UseCase, encrypter pkgCrt.Encrypter, secretKey string) Middleware {
	return Middleware{
		l:          l,
		jwtManager: jwtManager,
		userUC:     userUC,
		encrypter:  encrypter,
		secretKey:  secretKey,
	}
}
