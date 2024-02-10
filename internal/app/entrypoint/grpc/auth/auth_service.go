package auth

import (
	"github.com/vaberof/auth-grpc/internal/domain/auth"
	"github.com/vaberof/auth-grpc/pkg/domain"
)

type AuthService interface {
	Register(email domain.Email, password domain.Password) error
	Login(email domain.Email, password domain.Password) (*auth.AccessToken, error)
	Verify(email domain.Email, code domain.Code) error
	VerifyToken(token auth.AccessToken) error
}
