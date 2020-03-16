package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/abyanjksatu/goscription/models"
	"github.com/abyanjksatu/goscription/usecase"
	"github.com/abyanjksatu/goscription/util"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError will hold the response error structs
type ResponseError struct {
	Message string `json:"message"`
}

type articleHandler struct {
	AUsecase usecase.ArticleUsecase
}

// InitArticleHandler will initialize the article's HTTP handler
func InitArticleHandler(e *echo.Echo, us usecase.ArticleUsecase) {
	handler := &articleHandler{
		AUsecase: us,
	}
	e.GET("/articles", handler.FetchArticle)
	e.POST("/articles", handler.Store)
	e.GET("/articles/:id", handler.GetByID)
	e.DELETE("/articles/:id", handler.Delete)

}

// FetchArticle godoc
// @Summary Show a Article
// @Description get string by ID
// @Tags articles
// @Accept  json
// @Produce  json
// @Param num query string true "num"
// @Param cursor query string true "cursor"
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /articles [get]
func (a *articleHandler) FetchArticle(c echo.Context) error {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	listAr, nextCursor, err := a.AUsecase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listAr)
}

// FetchArticle godoc
// @Summary Show a Article
// @Description get string by ID
// @Tags articles
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "Article ID"
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /articles/{id} [get]
func (a *articleHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := int64(idP)
	art, err := a.AUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, art)
}

func isRequestValid(m *models.Article) (bool, error) {
	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store godoc
// @Summary Create an Article
// @Description Create an Article
// @Tags articles
// @Accept  json
// @Produce  json
// @Param article body models.Article true "Article Body"
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /articles [post]
func (a *articleHandler) Store(c echo.Context) error {
	var article models.Article
	err := c.Bind(&article)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&article); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.AUsecase.Store(ctx, &article)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, article)
}

// Delete godoc
// @Summary Create an Article
// @Description Create an Article
// @Tags articles
// @Accept  json
// @Produce  json
// @Param id query int true "Article ID"
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /articles [delete]
func (a *articleHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	id := int64(idP)
	err = a.AUsecase.Delete(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case util.ErrInternalServerError:
		return http.StatusInternalServerError
	case util.ErrNotFound:
		return http.StatusNotFound
	case util.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
