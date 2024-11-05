package v1

import (
	"fmt"

	"github.com/cutlery47/gostream/internal/service"
	"github.com/cutlery47/gostream/internal/storage"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var errMap = map[error]*echo.HTTPError{
	service.ErrChunkNotFound:         echo.ErrNotFound,
	service.ErrManifestNotFound:      echo.ErrNotFound,
	service.ErrVideoNotFound:         echo.ErrNotFound,
	service.ErrSegmentationException: echo.ErrInternalServerError,
	service.ErrNotImplemented:        echo.ErrNotImplemented,
	storage.ErrNotImplemented:        echo.ErrNotImplemented,
	storage.ErrUniueVideo:            echo.ErrBadRequest,
}

type errHandler struct {
	errLog *zap.Logger
}

func newErrHandler(errLog *zap.Logger) *errHandler {
	return &errHandler{
		errLog: errLog,
	}
}

func (h *errHandler) handle(err error) *echo.HTTPError {
	// trying to map err to HTTPError
	if httpErr, ok := errMap[err]; ok {
		httpErr.Message = err.Error()
		return httpErr
	}

	// log error if unexpected
	h.errLog.Error(fmt.Sprintf("Error: %v", err))

	// return 500 if:
	// 1) err couldn't be mapped to HTTPError
	// 2) err is an unexpected error
	return echo.ErrInternalServerError
}
