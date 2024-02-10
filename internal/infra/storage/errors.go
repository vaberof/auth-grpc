package storage

import "errors"

var (
	ErrPostgresUserNotFound = errors.New("user not found")

	ErrRedisKeyNotFound = errors.New("key not found")
)
