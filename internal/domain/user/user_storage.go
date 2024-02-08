package user

import (
	"context"
	"github.com/vaberof/auth-grpc/pkg/domain"
)

type UserStorage interface {
	Create(ctx context.Context, email domain.Email, password domain.Password) (domain.UserId, error)
	GetByEmail(ctx context.Context, email domain.Email) (*User, error)
	ExistsByEmail(ctx context.Context, email domain.Email) (bool, error)
}
