package mongo

import (
	"github.com/nguyenthanhtrung001/reup/internal/user/repository"
	"github.com/nguyenthanhtrung001/reup/pkg/jwt"
	"github.com/nguyenthanhtrung001/reup/pkg/log"
	"github.com/nguyenthanhtrung001/reup/pkg/mongo"
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
