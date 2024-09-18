package controller

import (
	"github.com/cutlery47/gostream/internal/service"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	echo   *echo.Echo
	router *router
}

func New(manifestService, chunkService service.Service) *Controller {
	e := echo.New()
	r := newRouter(manifestService, chunkService)

	e.GET("api/v1/:filename", r.demux)

	return &Controller{
		echo:   e,
		router: r,
	}
}

func (c *Controller) Handler() *echo.Echo {
	return c.echo
}
