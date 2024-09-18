package controller

import (
	"errors"
	"log"

	"github.com/labstack/echo/v4"
)

var (
	errInternalServer = errors.New("internal server error")
	errNotFound       = errors.New("resource not found")
)

type errHandler struct{}

func (eh errHandler) handle(err error) *echo.HTTPError {
	log.Println("An error occurred:", err)

	// error -> HTTPError mapping
	if errors.Is(err, errNotFound) {
		return echo.NewHTTPError(404, errNotFound)
	} else {
		return echo.NewHTTPError(500, errInternalServer)
	}
}
