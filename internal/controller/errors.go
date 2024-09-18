package controller

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var (
	errInternalServer = errors.New("internal server error")
	errNotFound       = errors.New("resource not found")
)

type errHandler struct {
	log *zap.Logger
}

func (eh errHandler) handle(err error) *echo.HTTPError {
	eh.log.Error(fmt.Sprintf("Error: %v", err))

	// error -> HTTPError mapping
	if errors.Is(err, errNotFound) {
		return echo.NewHTTPError(404, errNotFound)
	} else {
		return echo.NewHTTPError(500, errInternalServer)
	}
}
