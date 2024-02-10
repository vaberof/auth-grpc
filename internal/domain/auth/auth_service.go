package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vaberof/auth-grpc/internal/domain/user"
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

const (
	userDataCacheExpireTime         = 5 * time.Minute
	verificationCodeCacheExpireTime = 2 * time.Minute
)

const verificationCodeLength = 6

var (
	ErrUserAlreadyExists      = errors.New("user with specified email already exists")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")

	ErrInvalidVerificationCode = errors.New("invalid verification code")
	ErrVerificationCodeExpired = errors.New("verification code has expired")

	ErrTokenExpired = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type AuthService interface {
	Register(email domain.Email, password domain.Password) error
	Login(email domain.Email, password domain.Password) (*AccessToken, error)
	Verify(email domain.Email, code domain.Code) error
	VerifyToken(token AccessToken) error
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

func (a *authServiceImpl) Register(email domain.Email, password domain.Password) error {
	const operation = "Register"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.String("email", email.String()))

	log.Info("registering a user")

	exists, err := a.userService.ExistsByEmail(email)
	if err != nil {
		log.Error("failed to get info about existing/non-existing email")

		return fmt.Errorf("%s: %w", operation, err)
	}

	if exists {
		log.Warn("user already exists with specified email")

		return fmt.Errorf("%s: %w", operation, ErrUserAlreadyExists)
	}

	passwordHash, err := xpassword.Hash(password.String())
	if err != nil {
		log.Error("failed to hash password", "error", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	domainUser := &user.User{
		Email:    email,
		Password: domain.Password(passwordHash),
	}

	userData, err := json.Marshal(domainUser)
	if err != nil {
		log.Error("failed to marshal a user", "error", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	err = a.inMemoryStorage.Set(userEmailKey+email.String(), string(userData), userDataCacheExpireTime)
	if err != nil {
		log.Error("failed to set a userData to cache", "error", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	go func() {
		log.Info("send verification code to email")

		err := a.sendVerificationCode(registerCodeKey, email)
		if err != nil {
			log.Error("failed to send verification code", "error", err)
		}
	}()

	return nil
}

func (a *authServiceImpl) Login(email domain.Email, password domain.Password) (*AccessToken, error) {
	const operation = "Login"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.String("email", email.String()))

	log.Info("logging a user")

	domainUser, err := a.userService.GetByEmail(email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			log.Error("user not found", "error", err)

			return nil, fmt.Errorf("%s: %w", operation, ErrInvalidEmailOrPassword)
		}

		log.Error("failed to get user by email", "error", err)

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	err = xpassword.Check(password.String(), domainUser.Password.String())
	if err != nil {
		log.Error("incorrect password", err)

		return nil, fmt.Errorf("%s: %w", operation, ErrInvalidEmailOrPassword)
	}

	token, err := accesstoken.Create(domainUser.Id, a.config.TokenTtl, accesstoken.SecretKey(a.config.TokenSecretKey))
	if err != nil {
		log.Error("failed to create an access token", "error", err)

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	accessToken := AccessToken(token)

	log.Info("user logged in")

	return &accessToken, nil
}

func (a *authServiceImpl) Verify(email domain.Email, code domain.Code) error {
	const operation = "Verify"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.String("email", email.String()),
		slog.String("code", code.String()))

	log.Info("verifying an email")

	userData, err := a.inMemoryStorage.Get(userEmailKey + email.String())
	if err != nil {
		log.Error("failed to get user data from cache", "error", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	var domainUser user.User
	err = json.Unmarshal([]byte(userData), &domainUser)
	if err != nil {
		log.Error("failed to unmarshal user", "error", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	cachedCode, err := a.inMemoryStorage.Get(registerCodeKey + domainUser.Email.String())
	if err != nil {
		log.Error("failed to get verification code from cache", "error", err)

		return fmt.Errorf("%s: %w", operation, ErrVerificationCodeExpired)
	}

	if code.String() != cachedCode {
		log.Error("incorrect validation code", "error", ErrInvalidVerificationCode)

		return fmt.Errorf("%s: %w", operation, ErrInvalidVerificationCode)
	}

	_, err = a.userService.Create(domainUser.Email, domainUser.Password)
	if err != nil {
		log.Error("failed to create user", "error", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("user created")

	return nil
}

func (a *authServiceImpl) VerifyToken(token AccessToken) error {
	const operation = "VerifyToken"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.String("token", string(token)))

	log.Info("verifying a token")

	_, err := accesstoken.Verify(string(token), accesstoken.SecretKey(a.config.TokenSecretKey))
	if err != nil {
		if errors.Is(err, accesstoken.ErrInvalidToken) {
			log.Error("invalid access token", "error", err)

			return fmt.Errorf("%s: %w", operation, ErrInvalidToken)
		}
		if errors.Is(err, accesstoken.ErrExpiredToken) {
			log.Error("access token has expired", "error", err)

			return fmt.Errorf("%s: %w", operation, ErrTokenExpired)
		}
		if errors.Is(err, accesstoken.ErrInvalidSigningMethod) {
			log.Error("access token has invalid signing method", "error", err)

			return fmt.Errorf("%s: %w", operation, ErrInvalidToken)
		}

		log.Error("unexpected error from 'accesstoken' package", "error", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("successfully verified a token")

	return nil
}

func (a *authServiceImpl) sendVerificationCode(key string, email domain.Email) error {
	const operation = "sendVerificationCode"

	log := a.logger.With(
		slog.String("operation", operation),
		slog.String("key", key),
		slog.String("email", email.String()))

	code, err := xrand.GenerateRandomCode(verificationCodeLength)
	if err != nil {
		log.Error("failed to generate random code", "error", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("random code generated")

	err = a.inMemoryStorage.Set(key+email.String(), code, verificationCodeCacheExpireTime)
	if err != nil {
		log.Error("failed to cache verification code", "error", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("code saved in memory storage", "code", code)

	err = a.notificationService.SendEmail(email.String(), "verification_email", "Verification email", map[string]string{"code": code})
	if err != nil {
		log.Error("failed to send verification code", "error", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("code has been sent to notification server")

	return nil
}
