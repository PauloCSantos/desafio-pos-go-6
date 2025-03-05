package user_controller

import (
	"context"
	"fullcycle-auction_go/configuration/rest_err"
	"fullcycle-auction_go/internal/infra/api/web/validation"
	"fullcycle-auction_go/internal/usecase/user_usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *UserController) CreateUser(c *gin.Context) {
	var createUserInputoDTO user_usecase.CreateUserInputDTO

	if err := c.ShouldBindJSON(&createUserInputoDTO); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	id, err := u.userUseCase.CreateUser(context.Background(), createUserInputoDTO)

	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}
	c.JSON(http.StatusCreated, id)
	//c.Status(http.StatusCreated)
}
