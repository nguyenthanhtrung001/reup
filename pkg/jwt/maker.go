package jwt

import (
	"time"
)

// Token represents
//
//go:generate mockery --name Maker --inpackage
type Maker interface {
	CreateToken(payload Payload) (string, time.Time, error)

	VerifyToken(token string) (Payload, error)
}
