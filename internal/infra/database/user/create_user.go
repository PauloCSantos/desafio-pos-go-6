package user

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/user_entity"
	"fullcycle-auction_go/internal/internal_error"
	"github.com/google/uuid"
)

func (ur *UserRepository) CreateUser(
	ctx context.Context, userName string) (*user_entity.User, *internal_error.InternalError) {

	newUser := UserEntityMongo{
		Id:   uuid.New().String(),
		Name: userName,
	}

	_, err := ur.Collection.InsertOne(ctx, newUser)
	if err != nil {
		logger.Error("Error trying to create user", err)
		return nil, internal_error.NewInternalServerError("Error trying to create user")
	}

	userEntity := &user_entity.User{
		Id:   newUser.Id,
		Name: newUser.Name,
	}

	return userEntity, nil
}
