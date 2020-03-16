package http

import (
	"github.com/abyanjksatu/goscription/internal/http/middleware"
	"github.com/abyanjksatu/goscription/usecase"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Server is a struct
type Server struct {
	articleUsecase usecase.ArticleUsecase
	domainUsecase  usecase.DomainUsecase
}

// NewServer is constructor
func NewServer(du usecase.DomainUsecase, au usecase.ArticleUsecase) *Server {
	return &Server{
		domainUsecase:  du,
		articleUsecase: au,
	}
}

// Run function
func (s *Server) Run() {
	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)

	InitArticleHandler(e, s.articleUsecase)
	InitDomainHandler(e, s.domainUsecase)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Start(viper.GetString("server.address"))
}
