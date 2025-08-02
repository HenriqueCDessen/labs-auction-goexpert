package user_controller

import (
	"context"
	"fullcycle-auction_go/configuration/rest_err"
	"fullcycle-auction_go/internal/infra/api/web/validation"
	"fullcycle-auction_go/internal/usecase/user_usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUseCase user_usecase.UserUseCaseInterface
}

func NewUserController(userUseCase user_usecase.UserUseCaseInterface) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

func (u *UserController) CreateUser(c *gin.Context) {
	var userInputDTO user_usecase.UserOutputDTO

	if err := c.ShouldBindJSON(&userInputDTO); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	rest, err := u.userUseCase.CreateUser(context.Background(), userInputDTO.Name)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, rest)
}
