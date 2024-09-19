package controller

import (
	"fmt"

	"github.com/cutlery47/gostream/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type errHandler struct {
	log *zap.Logger
	// ServiceError -> echo.HTTPError mapping
	errMap map[*service.ServiceError]*echo.HTTPError
}

func newErrHandler(errLog *zap.Logger) *errHandler {
	errMap := map[*service.ServiceError]*echo.HTTPError{
		service.ErrManifestNotFound:      echo.ErrNotFound,
		service.ErrVideoNotFound:         echo.ErrNotFound,
		service.ErrSegmentationException: echo.ErrInternalServerError,
	}

	return &errHandler{
		log:    errLog,
		errMap: errMap,
	}
}

func (eh errHandler) handle(err error) *echo.HTTPError {
	eh.log.Error(fmt.Sprintf("Error: %v", err))

	// trying to map err to HTTPError
	if err, ok := err.(*service.ServiceError); ok {
		if httpErr, ok := eh.errMap[err]; ok {
			return httpErr
		}
	}

	// return 500 if:
	// 1) err couldn't be mapped to HTTPError
	// 2) err is an unexpected error
	return echo.ErrInternalServerError
}
