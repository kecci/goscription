package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/abyanjksatu/goscription/internal/service"
	"github.com/abyanjksatu/goscription/util"
	"github.com/labstack/echo/v4"
)

type userController struct {
	UService service.UserService
}

// InitUserController will initialize the article's HTTP controller
func InitUserController(e *echo.Echo, us service.UserService) {
	controller := &userController{
		UService: us,
	}
	// e.GET("/users", controller.FetchArticle)
	e.POST("/user", controller.Store)
	e.GET("/user/:id", controller.GetByID)
	// e.DELETE("/articles/:id", controller.Delete)
}

// GetUser godoc
// @Summary Show a User
// @Description get string by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /user/{id} [get]
func (a *userController) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(util.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := int64(idP)
	art, err := a.UService.GetByID(ctx, id)
	if err != nil {
		return c.JSON(util.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, art)
}

//UserRequest user body request
type UserRequest struct {
	Name   string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Store godoc
// @Summary Create an User
// @Description Create an User
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body UserRequest true "User Body"
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /user [post]
func (a *userController) Store(c echo.Context) error {
	var userRequest UserRequest
	err := c.Bind(&userRequest)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	userParam := service.UserParam{
		Name: userRequest.Name,
		Email: userRequest.Email,
		Password: userRequest.Password,
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	_, err = a.UService.Store(ctx, userParam)
	if err != nil {
		return c.JSON(util.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, userParam)
}