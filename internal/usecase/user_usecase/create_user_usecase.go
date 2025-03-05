package user_usecase

import (
	"context"
	"fullcycle-auction_go/internal/internal_error"
)

type CreateUserInputDTO struct {
	UserName string `json:"name"`
}

type CreateUserOutputDTO struct {
	Id string `json:"id"`
}

func (u *UserUseCase) CreateUser(
	ctx context.Context, input CreateUserInputDTO) (*CreateUserOutputDTO, *internal_error.InternalError) {
	userEntity, err := u.UserRepository.CreateUser(ctx, input.UserName)
	if err != nil {
		return nil, err
	}

	return &CreateUserOutputDTO{
		Id: userEntity.Id,
	}, nil
}
