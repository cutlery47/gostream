package controller

import (
	"fmt"

	"github.com/cutlery47/gostream/internal/service"
	"github.com/cutlery47/gostream/internal/storage"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type errHandler struct {
	log *zap.Logger
	// error -> echo.HTTPError mapping
	errMap map[error]*echo.HTTPError
}

func newErrHandler(errLog *zap.Logger) *errHandler {
	errMap := map[error]*echo.HTTPError{
		service.ErrChunkNotFound:         echo.ErrNotFound,
		service.ErrManifestNotFound:      echo.ErrNotFound,
		service.ErrVideoNotFound:         echo.ErrNotFound,
		service.ErrSegmentationException: echo.ErrInternalServerError,
		service.ErrNotImplemented:        echo.ErrNotImplemented,
		storage.ErrNotImplemented:        echo.ErrNotImplemented,
		storage.ErrUniueVideo:            echo.ErrBadRequest,
	}

	return &errHandler{
		log:    errLog,
		errMap: errMap,
	}
}

func (eh errHandler) handle(err error) *echo.HTTPError {
	// trying to map err to HTTPError
	if httpErr, ok := eh.errMap[err]; ok {
		httpErr.Message = err.Error()
		return httpErr
	}

	// log error if unexpected
	eh.log.Error(fmt.Sprintf("Error: %v", err))

	// return 500 if:
	// 1) err couldn't be mapped to HTTPError
	// 2) err is an unexpected error
	return echo.ErrInternalServerError
}
