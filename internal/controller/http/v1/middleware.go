package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func requestLoggerMiddleware(reqLog *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(
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
	)
}
