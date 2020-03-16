package http

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/abyanjksatu/goscription/usecase"
	"github.com/labstack/echo/v4"
)

type domainHandler struct {
	DUsecase usecase.DomainUsecase
}

// InitDomainHandler will initialize the domain's HTTP handler
func InitDomainHandler(e *echo.Echo, us usecase.DomainUsecase) {
	handler := &domainHandler{
		DUsecase: us,
	}
	e.GET("/domains", handler.GetDomains)
	e.GET("/domains/available", handler.GetDomainsAvailable)
}

// GetDomains godoc
// @Summary Show a Domains
// @Description get string by ID
// @Tags domains
// @Accept  json
// @Produce  json
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /domains [get]
func (d *domainHandler) GetDomains(c echo.Context) error {
	request, _ := http.NewRequest("GET", "https://api.ote-godaddy.com/v1/domains", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "sso-key 3mM44UaguLoR8V_S777TEwztnyJN8mQbAnGKD:7cMUQQQxaL3LpTpNoS9WqG")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}

	return c.JSON(http.StatusOK, response.Body)
}

// GetDomainsAvailable godoc
// @Summary Show a Domains Available
// @Description get string by ID
// @Tags domains
// @Accept  json
// @Produce  json
// @Param domain query string true "domain name"
// @Success 200 {object} models.DomainAvailableResponse
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /domains/available [get]
func (d *domainHandler) GetDomainsAvailable(c echo.Context) error {
	domain := c.QueryParam("domain")
	domainAvailableResponse, err := d.DUsecase.GetDomainAvailable(domain)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, domainAvailableResponse)
	}

	return c.JSON(http.StatusOK, domainAvailableResponse)
}
