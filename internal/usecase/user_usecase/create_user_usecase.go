package user_usecase

import (
	"context"
	"fullcycle-auction_go/internal/entity/user_entity"
	"fullcycle-auction_go/internal/internal_error"
)

type UserOutputDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserUseCaseInterface interface {
	CreateUser(
		ctx context.Context, name string) (*UserOutputDTO, *internal_error.InternalError)
	FindUserById(
		ctx context.Context,
		id string) (*UserOutputDTO, *internal_error.InternalError)
}

func NewUserUseCase(userRepository user_entity.UserRepositoryInterface) UserUseCaseInterface {
	return &UserUseCase{
		userRepository,
	}
}

type UserUseCase struct {
	UserRepository user_entity.UserRepositoryInterface
}

func (u *UserUseCase) CreateUser(
	ctx context.Context, name string) (*UserOutputDTO, *internal_error.InternalError) {

	userEntity, err := user_entity.CreateUser(name)
	if err != nil {
		return nil, err
	}

	if repoErr := u.UserRepository.CreateUser(ctx, *userEntity); repoErr != nil {
		return nil, repoErr
	}

	return &UserOutputDTO{
		Id:   userEntity.Id,
		Name: userEntity.Name,
	}, nil
}
