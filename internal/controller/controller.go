package controller

import (
	"github.com/cutlery47/gostream/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Controller struct {
	echo   *echo.Echo
	router *router
	log    *zap.Logger
}

func New(
	service service.Service,
	reqLog, errLog, infoLog *zap.Logger) *Controller {
	e := echo.New()
	r := newRouter(errLog, service)

	e.Use(
		middleware.RequestLoggerWithConfig(
			middleware.RequestLoggerConfig{
				LogMethod:   true,
				LogStatus:   true,
				LogRemoteIP: true,
				LogURI:      true,
				LogError:    true,
				LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
					reqLog.Info(
						"",
						zap.String("method", v.Method),
						zap.Int("status", v.Status),
						zap.String("IP", v.RemoteIP),
						zap.String("URI", v.URI),
						zap.Error(v.Error),
					)
					return nil
				},
			},
		),
	)

	// file retrieval
	e.GET("api/v1/:filename", r.getFile)

	// file upload
	e.POST("api/v1/upload", r.uploadFile)

	// file deletion
	e.DELETE("api/v1/:filename", r.deleteFile)

	return &Controller{
		echo:   e,
		router: r,
		log:    infoLog,
	}
}

func (c *Controller) Handler() *echo.Echo {
	return c.echo
}
