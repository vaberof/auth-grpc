package auth

import (
	"context"
	"errors"
	"github.com/vaberof/auth-grpc/pkg/domain"
	"github.com/vaberof/auth-grpc/pkg/logging/logs"
	"time"
)

var (
	ErrUserAlreadyExists      = errors.New("user with specified email already exists")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")

	ErrTokenExpired = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type AuthService interface {
	RegisterUser(ctx context.Context, email domain.Email, password domain.Password) (domain.UserId, error)
	Login(ctx context.Context, email domain.Email, password domain.Password, appId int32) (*AccessToken, error)
	VerifyToken(ctx context.Context, token *AccessToken) error
}

type Config struct {
	TokenTtl       time.Duration `yaml:"token-ttl"`
	TokenSecretKey string        `yaml:"access-token-secret-key"`
}

type authServiceImpl struct {
	config          *Config
	userService     UserService
	inMemoryStorage InMemoryStorage

	logger *logs.Logs
}

func NewAuthService(config *Config, userService UserService, inMemoryStorage InMemoryStorage, logger *logs.Logs) AuthService {
	return &authServiceImpl{config: config, userService: userService, inMemoryStorage: inMemoryStorage, logger: logger}
}

func (service *authServiceImpl) RegisterUser(ctx context.Context, email domain.Email, password domain.Password) (domain.UserId, error) {
	//TODO implement me
	panic("implement me")
}

func (service *authServiceImpl) Login(ctx context.Context, email domain.Email, password domain.Password, appId int32) (*AccessToken, error) {
	//TODO implement me
	panic("implement me")
}

func (service *authServiceImpl) VerifyToken(ctx context.Context, token *AccessToken) error {
	//TODO implement me
	panic("implement me")
}
