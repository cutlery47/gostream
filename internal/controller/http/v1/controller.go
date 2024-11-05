package v1

import (
	"github.com/cutlery47/gostream/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

func NewController(e *echo.Echo, s service.Service, reqLog, errLog, infoLog *zap.Logger) {
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error { return c.NoContent(200) })
	e.GET("/swagger", echoSwagger.WrapHandler)

	v1 := e.Group("/api/v1", requestLoggerMiddleware(reqLog))
	{
		newFileRoutes(v1.Group("/files"), s, newErrHandler(errLog))
	}
}
