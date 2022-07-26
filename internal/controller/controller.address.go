package controller

import (
	"fmt"
	"net/http"

	"github.com/kecci/goscription/internal/service"
	"github.com/kecci/goscription/models"
	"github.com/labstack/echo/v4"
)

type addressController struct {
	addressService service.AddressService
}

func NewAddressController(e *echo.Echo, addressService service.AddressService) {
	controller := &addressController{
		addressService: addressService,
	}

	e.GET("/address", controller.GetAddressAll)
	e.POST("/address", controller.Insert)
}

// GetAddressAll godoc
// @Summary Show a Address
// @Description get address
// @Tags address
// @Accept json
// @Produce json
// @Header 200 {string} models.BaseResponse
// @Failure 400 {object} models.BaseResponse
// @Failure 404 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /address [get]
func (a *addressController) GetAddressAll(c echo.Context) error {
	addresses, err := a.addressService.GetAddressAll()
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, models.BaseResponse{Code: "FAILED", Message: "FAILED", Error: []string{err.Error()}})
	}
	return c.JSON(http.StatusOK, models.BaseResponse{Code: "SUCCESS", Message: "SUCCESS", Data: addresses})
}

// Insert godoc
// @Summary Show a Address
// @Description insert address
// @Tags address
// @Accept json
// @Produce json
// @Param address body models.Address true "address"
// @Header 200 {string} models.BaseResponse
// @Failure 400 {object} models.BaseResponse
// @Failure 404 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /address [post]
func (a *addressController) Insert(c echo.Context) error {
	var address models.Address
	if err := c.Bind(&address); err != nil {
		return c.JSON(http.StatusBadRequest, models.BaseResponse{Code: "FAILED", Message: "FAILED", Error: []string{err.Error()}})
	}
	err := a.addressService.Insert(address)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, models.BaseResponse{Code: "FAILED", Message: "FAILED", Error: []string{err.Error()}})
	}
	return c.JSON(http.StatusOK, models.BaseResponse{Code: "SUCCESS", Message: "SUCCESS", Data: true})
}
