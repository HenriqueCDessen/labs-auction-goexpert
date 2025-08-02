package user_entity

import (
	"context"
	"fullcycle-auction_go/internal/internal_error"

	"github.com/google/uuid"
)

type User struct {
	Id   string
	Name string
}

func CreateUser(name string) (*User, *internal_error.InternalError) {
	user := &User{
		Id:   uuid.New().String(),
		Name: name,
	}

	return user, nil
}

type UserRepositoryInterface interface {
	CreateUser(
		ctx context.Context, user User) *internal_error.InternalError
	FindUserById(
		ctx context.Context, userId string) (*User, *internal_error.InternalError)
}
