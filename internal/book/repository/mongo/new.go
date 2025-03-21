package mogo

import (
	"reup/internal/book/repository"
	"reup/pkg/jwt"
	"reup/pkg/log"
	"reup/pkg/mongo"
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
