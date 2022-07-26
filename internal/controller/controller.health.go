package controller

import (
	"context"
	"net/http"

	"github.com/kecci/goscription/internal/service"
	"github.com/kecci/goscription/models"
	"github.com/kecci/goscription/utility"
	"github.com/labstack/echo/v4"
)

type healthController struct {
	healthService service.HealthService
}

// InitHealthController will initialize the health's HTTP controller
func InitHealthController(e *echo.Echo, healthService service.HealthService) {
	controller := &healthController{
		healthService: healthService,
	}
	e.GET("/health", controller.CheckHealth)
}

// CheckHealth godoc
// @Summary Show a Health
// @Description get health
// @Tags health
// @Accept json
// @Produce json
// @Header 200 {string} models.BaseResponse
// @Failure 400 {object} models.BaseResponse
// @Failure 404 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /health [get]
func (h healthController) CheckHealth(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	err := h.healthService.CheckHealth(ctx)
	if err != nil {
		return c.JSON(utility.GetStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, models.BaseResponse{Code: "SUCCESS", Message: "SUCCESS", Data: true})
}
