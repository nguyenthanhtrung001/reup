package jwt

import (
	"errors"
	"time"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	UserID    string    `json:"user_id"`
	GroupID   string    `json:"group_id"`
	GroupRole string    `json:"group_role"`
	ExpiredAt time.Time `json:"expired_at"`
	Exp       int64     `json:"exp"`
}

type UserField struct {
	UserID    string
	GroupID   string
	GroupRole string
}

func NewPayload(profile UserField, duration time.Duration) Payload {
	exp := time.Now().Add(duration)

	return Payload{
		UserID:    profile.UserID,
		GroupID:   profile.GroupID,
		GroupRole: profile.GroupRole,
		ExpiredAt: exp,
		Exp:       exp.UTC().Unix(),
	}
}

func (p Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
