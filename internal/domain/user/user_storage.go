package user

import (
	"github.com/vaberof/auth-grpc/pkg/domain"
)

type UserStorage interface {
	Create(email domain.Email, password domain.Password) (domain.UserId, error)
	GetByEmail(email domain.Email) (*User, error)
	ExistsByEmail(email domain.Email) (bool, error)
}
