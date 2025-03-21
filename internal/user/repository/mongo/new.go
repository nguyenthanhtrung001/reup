package mongo

import (
	"reup/internal/user/repository"
	"reup/pkg/jwt"
	"reup/pkg/log"
	"reup/pkg/mongo"
)

type implRepository struct {
	l          log.Logger
	database   mongo.Database
	jwtManager jwt.Maker
}

var _ repository.Repository = implRepository{}

func New(l log.Logger, database mongo.Database, jwtManager jwt.Maker) repository.Repository {
	return implRepository{
		l:          l,
		database:   database,
		jwtManager: jwtManager,
	}
}
