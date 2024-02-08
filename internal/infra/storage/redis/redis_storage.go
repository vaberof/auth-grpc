package redis

import "errors"

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExpired  = errors.New("key has expired")
)
