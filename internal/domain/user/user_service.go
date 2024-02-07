package user

import (
	"context"
	"github.com/vaberof/auth-grpc/pkg/domain"
	"github.com/vaberof/auth-grpc/pkg/logging/logs"
)

var (
	ErrUserNotFound = "user with specified email not found"
)

type UserService interface {
	Create(ctx context.Context, email domain.Email, password domain.Password) (domain.UserId, error)
	GetByEmail(ctx context.Context, email domain.Email) (*User, error)
	ExistsByEmail(ctx context.Context, email domain.Email) bool
}

type userServiceImpl struct {
	userStorage UserStorage

	logger *logs.Logs
}

func New(userStorage UserStorage, logger *logs.Logs) UserService {
	return &userServiceImpl{userStorage: userStorage, logger: logger}
}

func (service *userServiceImpl) Create(ctx context.Context, email domain.Email, password domain.Password) (domain.UserId, error) {
	//TODO implement me
	panic("implement me")
}

func (service *userServiceImpl) GetByEmail(ctx context.Context, email domain.Email) (*User, error) {
	//TODO implement me
	panic("implement me")
}

func (service *userServiceImpl) ExistsByEmail(ctx context.Context, email domain.Email) bool {
	//TODO implement me
	panic("implement me")
}
