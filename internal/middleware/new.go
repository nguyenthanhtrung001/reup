package middleware

import (
	userUseCase "reup/internal/user/usecase"
	pkgCrt "reup/pkg/encrypter"
	"reup/pkg/jwt"
	"reup/pkg/log"
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
