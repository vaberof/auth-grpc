package auth

import (
	"github.com/vaberof/auth-grpc/pkg/domain"
	"time"
)

type JwtPayload struct {
	UserId    domain.UserId
	ExpiredAt time.Time
}

func NewPayload(userId domain.UserId, ttl time.Duration) *JwtPayload {
	return &JwtPayload{
		UserId:    userId,
		ExpiredAt: time.Now().Add(ttl),
	}
}
