package user

import "github.com/vaberof/auth-grpc/pkg/domain"

type User struct {
	Id       domain.UserId
	Email    domain.Email
	Password domain.Password
}
