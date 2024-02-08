package accesstoken

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vaberof/auth-grpc/pkg/auth"
	"github.com/vaberof/auth-grpc/pkg/domain"
	"strconv"
	"time"
)

var (
	ErrInvalidToken         = errors.New("token is invalid")
	ErrInvalidSigningMethod = errors.New("signing method is invalid")
	ErrExpiredToken         = errors.New("token has expired")
)

type SecretKey string

// Create returns JWT-token signed with specified secret key and
// stores UserId, ExpireAt and IssuedAt in jwt payload
func Create(userId domain.UserId, ttl time.Duration, secretKey SecretKey) (string, error) {
	payload := auth.NewPayload(userId, ttl)

	jwtWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    payload.UserId.String(),
		IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
		ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
	})

	token, err := jwtWithClaims.SignedString([]byte(secretKey))

	return token, err
}

// CreateWithExpirationTime is the same as Create, but additionally returns token expiration
func CreateWithExpirationTime(userId domain.UserId, ttl time.Duration, secretKey SecretKey) (string, time.Time, error) {
	payload := auth.NewPayload(userId, ttl)

	jwtWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    payload.UserId.String(),
		IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
		ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
	})

	token, err := jwtWithClaims.SignedString([]byte(secretKey))

	return token, payload.ExpiredAt, err
}

func Verify(token string, secretKey SecretKey) (*auth.JwtPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	if hasExpired(claims.ExpiresAt.Time) {
		return nil, ErrExpiredToken
	}

	uid, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		return nil, err
	}

	payload := &auth.JwtPayload{
		UserId:    domain.UserId(uid),
		IssuedAt:  claims.IssuedAt.Time,
		ExpiredAt: claims.ExpiresAt.Time,
	}

	return payload, nil
}

func hasExpired(expireTime time.Time) bool {
	currentTime := time.Now().UTC()
	return currentTime.After(expireTime.UTC())
}
