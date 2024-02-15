package user

import (
	"errors"
	"fmt"
	"github.com/vaberof/auth-grpc/internal/infra/storage"
	"github.com/vaberof/auth-grpc/pkg/domain"
	"github.com/vaberof/auth-grpc/pkg/logging/logs"
	"log/slog"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	Create(email domain.Email, password domain.Password) (domain.UserId, error)
	GetByEmail(email domain.Email) (*User, error)
	ExistsByEmail(email domain.Email) (bool, error)
}

type userServiceImpl struct {
	userStorage UserStorage

	logger *slog.Logger
}

func NewUserService(userStorage UserStorage, logs *logs.Logs) UserService {
	logger := logs.WithName("domain.user.service")
	return &userServiceImpl{userStorage: userStorage, logger: logger}
}

func (u *userServiceImpl) Create(email domain.Email, password domain.Password) (domain.UserId, error) {
	const operation = "Create"

	log := u.logger.With(
		slog.String("operation", operation),
		slog.String("email", string(email)))

	log.Info("creating a user")

	uid, err := u.userStorage.Create(email, password)
	if err != nil {
		log.Error("failed to create a user", "error", err)

		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("user created")

	return uid, nil
}

func (u *userServiceImpl) GetByEmail(email domain.Email) (*User, error) {
	const operation = "GetByEmail"

	log := u.logger.With(
		slog.String("operation", operation),
		slog.String("email", string(email)))

	domainUser, err := u.userStorage.GetByEmail(email)
	if err != nil {
		if errors.Is(err, storage.ErrPostgresUserNotFound) {
			log.Error("user with given email not found", "error", err)

			return nil, fmt.Errorf("%s: %w", operation, err)
		}

		log.Error("unexpected error from user storage", "error", err)

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("received user by email")

	return domainUser, nil
}

func (u *userServiceImpl) ExistsByEmail(email domain.Email) (bool, error) {
	const operation = "GetByEmail"

	log := u.logger.With(
		slog.String("operation", operation),
		slog.String("email", string(email)))

	exists, err := u.userStorage.ExistsByEmail(email)
	if err != nil {
		log.Error("failed to get info about existing/non-existing email", "error", err)

		return false, fmt.Errorf("%s: %w", operation, err)
	}

	return exists, nil
}
