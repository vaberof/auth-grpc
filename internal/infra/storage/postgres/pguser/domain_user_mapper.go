package pguser

import (
	"github.com/vaberof/auth-grpc/internal/domain/user"
	"github.com/vaberof/auth-grpc/pkg/domain"
)

func toDomainUser(pgUser *User) *user.User {
	return &user.User{
		Id:       domain.UserId(pgUser.Id),
		Email:    domain.Email(pgUser.Email),
		Password: domain.Password(pgUser.Password),
	}
}
