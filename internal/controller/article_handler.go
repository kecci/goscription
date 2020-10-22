package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/abyanjksatu/goscription/internal/service"
	"github.com/abyanjksatu/goscription/util"
	"github.com/labstack/echo/v4"
)

type articleController struct {
	AService service.ArticleService
}

// InitArticleController will initialize the article's HTTP controller
func InitArticleController(e *echo.Echo, us service.ArticleService) {
	controller := &articleController{
		AService: us,
	}
	e.GET("/articles", controller.FetchArticle)
	e.POST("/articles", controller.Store)
	e.GET("/articles/:id", controller.GetByID)
	e.DELETE("/articles/:id", controller.Delete)
}

// ResponseError will hold the response error structs
type ResponseError struct {
	Message string `json:"message"`
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
func (a *articleController) FetchArticle(c echo.Context) error {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	listAr, nextCursor, err := a.AService.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(util.GetStatusCode(err), ResponseError{Message: err.Error()})
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
func (a *articleController) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(util.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := int64(idP)
	art, err := a.AService.GetByID(ctx, id)
	if err != nil {
		return c.JSON(util.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, art)
}

//ArticleRequest article body request
type ArticleRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

// Store godoc
// @Summary Create an Article
// @Description Create an Article
// @Tags articles
// @Accept  json
// @Produce  json
// @Param article body ArticleRequest true "Article Body"
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /articles [post]
func (a *articleController) Store(c echo.Context) error {
	var articleRequest ArticleRequest
	err := c.Bind(&articleRequest)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	articleParam := service.ArticleParam{
		Title:   articleRequest.Title,
		Content: articleRequest.Content,
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.AService.Store(ctx, articleParam)
	if err != nil {
		return c.JSON(util.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, articleParam)
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
func (a *articleController) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(util.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	id := int64(idP)
	err = a.AService.Delete(ctx, id)
	if err != nil {
		return c.JSON(util.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Update godoc
// @Summary Update an Article
// @Description Update an Article
// @Tags articles
// @Accept  json
// @Produce  json
// @Param article body ArticleRequest true "Article Body"
// @Param id path int true "Article ID"
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /articles/{id} [put]
func (a *articleController) Update(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(util.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	var ar ArticleRequest
	err = c.Bind(&ar)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	articleParam := service.ArticleParam{
		ID: int64(idP),
		Title:   ar.Title,
		Content: ar.Content,
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.AService.Update(ctx, articleParam)
	if err != nil {
		return c.JSON(util.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, articleParam)
}