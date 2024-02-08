package auth

import (
	"github.com/vaberof/auth-grpc/pkg/domain"
	"time"
)

type JwtPayload struct {
	UserId    domain.UserId
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func NewPayload(userId domain.UserId, ttl time.Duration) *JwtPayload {
	return &JwtPayload{
		UserId:    userId,
		IssuedAt:  time.Now().UTC(),
		ExpiredAt: time.Now().UTC().Add(ttl),
	}
}
