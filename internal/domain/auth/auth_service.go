package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vaberof/auth-grpc/internal/domain/user"
	"github.com/vaberof/auth-grpc/internal/infra/storage/redis"
	"github.com/vaberof/auth-grpc/pkg/auth/accesstoken"
	"github.com/vaberof/auth-grpc/pkg/domain"
	"github.com/vaberof/auth-grpc/pkg/logging/logs"
	"github.com/vaberof/auth-grpc/pkg/xpassword"
	"github.com/vaberof/auth-grpc/pkg/xrand"
	"log/slog"
	"time"
)

const (
	userEmailKey    = "user_email_"
	registerCodeKey = "register_code_"
)

var (
	ErrUserAlreadyExists      = errors.New("user with specified email already exists")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")

	ErrInvalidVerificationCode = errors.New("invalid verification code")
	ErrVerificationCodeExpired = errors.New("verification code has expired")

	ErrTokenExpired = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type AuthService interface {
	Register(ctx context.Context, email domain.Email, password domain.Password) error
	Login(ctx context.Context, email domain.Email, password domain.Password) (*AccessToken, error)
	Verify(ctx context.Context, email domain.Email, code domain.Code) error
	VerifyToken(ctx context.Context, token AccessToken) error
}

type Config struct {
	TokenTtl       time.Duration `yaml:"token-ttl"`
	TokenSecretKey string        `yaml:"token-secret-key"`
}

type authServiceImpl struct {
	config              *Config
	userService         UserService
	notificationService NotificationService
	inMemoryStorage     InMemoryStorage

	logger *slog.Logger
}

func NewAuthService(config *Config, userService UserService, notificationService NotificationService, inMemoryStorage InMemoryStorage, logs *logs.Logs) AuthService {
	logger := logs.WithName("domain.auth.service")
	return &authServiceImpl{
		config:              config,
		userService:         userService,
		notificationService: notificationService,
		inMemoryStorage:     inMemoryStorage,
		logger:              logger,
	}
}

func (a *authServiceImpl) Register(ctx context.Context, email domain.Email, password domain.Password) error {
	const operation = "Register"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.String("email", email.String()))

	log.Info("registering a user")

	exists, err := a.userService.ExistsByEmail(ctx, email)
	if err != nil {
		log.Error("failed to get info about existing/non-existing email")

		return fmt.Errorf("%s: %w", operation, err)
	}

	if exists {
		log.Warn("user already exists with specified email", email)

		return ErrUserAlreadyExists
	}

	passwordHash, err := xpassword.Hash(password.String())
	if err != nil {
		log.Error("failed to hash password", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	domainUser := &user.User{
		Email:    email,
		Password: domain.Password(passwordHash),
	}

	userData, err := json.Marshal(domainUser)
	if err != nil {
		log.Error("failed to marshal a user", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	err = a.inMemoryStorage.Set(ctx, userEmailKey+email.String(), string(userData), 10*time.Minute)
	if err != nil {
		log.Error("failed to set a userData to cache", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	go func() {
		log.Info("send verification code to email", email)

		err := a.sendVerificationCode(ctx, registerCodeKey, email)
		if err != nil {
			log.Error("failed to send verification code", err)
		}
	}()

	return nil
}

func (a *authServiceImpl) Login(ctx context.Context, email domain.Email, password domain.Password) (*AccessToken, error) {
	const operation = "Login"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.String("email", email.String()))

	log.Info("logging a user")

	domainUser, err := a.userService.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			log.Error("user not found", err)

			return nil, fmt.Errorf("%s: %w", operation, ErrInvalidEmailOrPassword)
		}

		log.Error("failed to get user by email", err)

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	err = xpassword.Check(string(password), string(domainUser.Password))
	if err != nil {
		log.Error("incorrect password", err)

		return nil, fmt.Errorf("%s: %w", operation, ErrInvalidEmailOrPassword)
	}

	token, expiredAt, err := accesstoken.CreateWithExpirationTime(domainUser.Id, a.config.TokenTtl, accesstoken.SecretKey(a.config.TokenSecretKey))
	if err != nil {
		log.Error("failed to create an access token", err)

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	accessToken := AccessToken(token)

	err = a.inMemoryStorage.Set(ctx, token, domainUser.Id.String(), expiredAt.Sub(time.Now().UTC()))
	if err != nil {
		log.Error("failed to cache token", err)

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("user logged in")

	return &accessToken, nil
}

func (a *authServiceImpl) Verify(ctx context.Context, email domain.Email, code domain.Code) error {
	const operation = "Verify"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.String("email", email.String()),
		slog.String("code", code.String()))

	log.Info("verifying an email")

	userData, err := a.inMemoryStorage.Get(ctx, userEmailKey+email.String())
	if err != nil {
		log.Error("failed to get user data from cache", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	var domainUser user.User
	err = json.Unmarshal([]byte(userData), &domainUser)
	if err != nil {
		log.Error("failed to unmarshal user", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	cachedCode, err := a.inMemoryStorage.Get(ctx, registerCodeKey+domainUser.Email.String())
	if err != nil {
		log.Error("failed to get verification code from cache", err)

		return fmt.Errorf("%s: %w", operation, ErrVerificationCodeExpired)
	}

	if code.String() != cachedCode {
		log.Error("incorrect validation code", ErrInvalidVerificationCode)

		return fmt.Errorf("%s: %w", operation, ErrInvalidVerificationCode)
	}

	_, err = a.userService.Create(ctx, domainUser.Email, domainUser.Password)
	if err != nil {
		log.Error("failed to create user", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("user created")

	return nil
}

func (a *authServiceImpl) VerifyToken(ctx context.Context, token AccessToken) error {
	const operation = "VerifyToken"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.String("token", string(token)))

	log.Info("verifying a token")

	_, err := a.inMemoryStorage.Get(ctx, string(token))
	if err != nil {
		if errors.Is(err, redis.ErrKeyNotFound) {
			log.Error("specified token not found", err)

			return fmt.Errorf("%s: %w", operation, ErrInvalidToken)
		}

		if errors.Is(err, redis.ErrKeyExpired) {
			log.Error("specified token has expired", err)

			return fmt.Errorf("%s: %w", operation, ErrTokenExpired)
		}

		log.Error("unexpected error from memory storage", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

func (a *authServiceImpl) sendVerificationCode(ctx context.Context, key string, email domain.Email) error {
	const operation = "sendVerificationCode"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.String("key", key),
		slog.String("email", email.String()))

	code, err := xrand.GenerateRandomCode(6)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("random code generated")

	err = a.inMemoryStorage.Set(ctx, key+email.String(), code, time.Minute)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("code saved in memory storage")

	err = a.notificationService.SendEmail(ctx, email.String(), "verification_email", "Verification email", map[string]string{"code": code})
	if err != nil {
		log.Error("failed to send verification code", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("code has been sent to notification server")

	return nil
}
