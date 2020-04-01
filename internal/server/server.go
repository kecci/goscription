package server

import (
	"context"

	"github.com/abyanjksatu/goscription/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"
)

//Module for server
var Module = fx.Provide(NewServer)

//NewServer initialize new server
func NewServer(lc fx.Lifecycle) *echo.Echo {
	instance := echo.New()

	// Middleware
	middL := middleware.InitMiddleware()
	instance.Use(middL.CORS)
	instance.Use(middL.Logger)
	instance.Use(middL.Recover)

	instance.HTTPErrorHandler = middL.ErrorHandler

	instance.GET("/swagger/*", echoSwagger.WrapHandler)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logrus.Print("Starting HTTP server.")
			go instance.Start(viper.GetString("server.address"))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logrus.Print("Stopping HTTP server.")
			return instance.Shutdown(ctx)
		},
	})
	return instance
}
