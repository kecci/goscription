package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

//GoMiddleware struct of middleware
type GoMiddleware struct {
	// another stuff , may be needed by middleware
}

//CORS set cors by echo
func (m *GoMiddleware) CORS(h echo.HandlerFunc) echo.HandlerFunc {
	cors := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	})

	return cors(h)
}

//Recover set recover by echo
func (m *GoMiddleware) Recover(h echo.HandlerFunc) echo.HandlerFunc {
	recover := middleware.Recover()
	return recover(h)
}

func makeLogEntry(c echo.Context) *log.Entry {
	if c == nil {
		return log.WithFields(log.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	return log.WithFields(log.Fields{
		"at":     time.Now().Format("2006-01-02 15:04:05"),
		"method": c.Request().Method,
		"uri":    c.Request().URL.String(),
		"ip":     c.Request().RemoteAddr,
	})
}

//Logger for log
func (m *GoMiddleware) Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		makeLogEntry(c).Info("incoming request")
		return next(c)
	}
}

//ErrorHandler fot context echo
func (m *GoMiddleware) ErrorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if ok {
		report.Message = fmt.Sprintf("http error %d - %v", report.Code, report.Message)
	} else {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	makeLogEntry(c).Error(report.Message)
	c.HTML(report.Code, report.Message.(string))
}

//InitMiddleware will initialize the middleware handler
func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}
