package auth

import (
	"context"
	"github.com/vaberof/auth-grpc/internal/domain/auth"
	"github.com/vaberof/auth-grpc/pkg/domain"
)

type AuthService interface {
	Register(ctx context.Context, email domain.Email, password domain.Password) error
	Login(ctx context.Context, email domain.Email, password domain.Password) (*auth.AccessToken, error)
	Verify(ctx context.Context, email domain.Email, code domain.Code) error
	VerifyToken(ctx context.Context, token auth.AccessToken) error
}
