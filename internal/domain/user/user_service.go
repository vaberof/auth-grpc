package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/vaberof/auth-grpc/internal/infra/storage/postgres/pguser"
	"github.com/vaberof/auth-grpc/pkg/domain"
	"github.com/vaberof/auth-grpc/pkg/logging/logs"
	"log/slog"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	Create(ctx context.Context, email domain.Email, password domain.Password) (domain.UserId, error)
	GetByEmail(ctx context.Context, email domain.Email) (*User, error)
	ExistsByEmail(ctx context.Context, email domain.Email) (bool, error)
}

type userServiceImpl struct {
	userStorage UserStorage

	logger *slog.Logger
}

func NewUserService(userStorage UserStorage, logs *logs.Logs) UserService {
	logger := logs.WithName("domain.user.service")
	return &userServiceImpl{userStorage: userStorage, logger: logger}
}

func (u *userServiceImpl) Create(ctx context.Context, email domain.Email, password domain.Password) (domain.UserId, error) {
	const operation = "Create"

	log := u.logger.With(
		slog.String("operation", operation),
		slog.String("email", string(email)))

	log.Info("creating a user")

	return 1, nil

	// TODO: implement storage

	//uid, err := u.userStorage.Create(ctx, email, password)
	//if err != nil {
	//	log.Error("failed to create a user", err)
	//
	//	return 0, fmt.Errorf("%s: %w", operation, err)
	//}
	//
	//log.Info("used created")
	//
	//return uid, nil
}

func (u *userServiceImpl) GetByEmail(ctx context.Context, email domain.Email) (*User, error) {
	const operation = "GetByEmail"

	log := u.logger.With(
		slog.String("operation", operation),
		slog.String("email", string(email)))

	domainUser, err := u.userStorage.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pguser.ErrUserNotFound) {
			log.Error("user with given email not found", err)

			return nil, fmt.Errorf("%s: %w", operation, err)
		}

		log.Error("unexpected error from user storage", err)

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("received user by email")

	return domainUser, nil
}

func (u *userServiceImpl) ExistsByEmail(ctx context.Context, email domain.Email) (bool, error) {
	const operation = "GetByEmail"

	log := u.logger.With(
		slog.String("operation", operation),
		slog.String("email", string(email)))

	exists, err := u.ExistsByEmail(ctx, email)
	if err != nil {
		log.Error("failed to get info about existing/non-existing email", err)

		return false, fmt.Errorf("%s: %w", operation, err)
	}

	return exists, nil
}
