package auth

import (
	"github.com/vaberof/auth-grpc/internal/domain/user"
	"github.com/vaberof/auth-grpc/pkg/domain"
)

type UserService interface {
	Create(email domain.Email, password domain.Password) (domain.UserId, error)
	GetByEmail(email domain.Email) (*user.User, error)
	ExistsByEmail(email domain.Email) (bool, error)
}
