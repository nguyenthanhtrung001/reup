package mogo

import (
	"book-store/internal/book/repository"
	"book-store/pkg/jwt"
	"book-store/pkg/log"
	"book-store/pkg/mongo"
)

type implRepository struct {
	l          log.Logger
	database   mongo.Database
	jwtManager jwt.Maker
}

var _ repository.Repository = implRepository{} // dấu '_' bắt buộc phải impl

func New(l log.Logger, database mongo.Database, jwtManager jwt.Maker) repository.Repository {
	return implRepository{
		l:          l,
		database:   database,
		jwtManager: jwtManager,
	}
}
