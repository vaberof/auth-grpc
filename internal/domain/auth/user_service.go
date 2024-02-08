package auth

import (
	"context"
	"github.com/vaberof/auth-grpc/internal/domain/user"
	"github.com/vaberof/auth-grpc/pkg/domain"
)

type UserService interface {
	Create(ctx context.Context, email domain.Email, password domain.Password) (domain.UserId, error)
	GetByEmail(ctx context.Context, email domain.Email) (*user.User, error)
	ExistsByEmail(ctx context.Context, email domain.Email) (bool, error)
}
