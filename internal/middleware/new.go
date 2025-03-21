package middleware

import (
	userUseCase "book-store/internal/user/usecase"
	pkgCrt "book-store/pkg/encrypter"
	"book-store/pkg/jwt"
	"book-store/pkg/log"
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
