package mogo

import (
	"github.com/nguyenthanhtrung001/reup/internal/scan_video/repository"
	"github.com/nguyenthanhtrung001/reup/pkg/jwt"
	"github.com/nguyenthanhtrung001/reup/pkg/log"
	"github.com/nguyenthanhtrung001/reup/pkg/mongo"
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
