package controller

import (
	"io"
	"mime/multipart"
	"strings"

	"github.com/cutlery47/gostream/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type router struct {
	service    service.Service
	errHandler errHandler
}

func newRouter(errLog *zap.Logger, service service.Service) *router {
	return &router{
		service:    service,
		errHandler: *newErrHandler(errLog),
	}
}

// GET /api/v1/:filename
func (r *router) getFile(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	filename := c.Param("filename")

	return r.get(c, filename)
}

func (r *router) get(c echo.Context, filename string) (err error) {
	var file io.ReadCloser

	ctx := c.Request().Context()

	// searching for requested file
	file, err = r.service.Serve(ctx, filename)
	if err != nil {
		return r.errHandler.handle(err)
	}

	// converting the file into a sequence of bytes
	blob, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// returning the file
	return c.Blob(200, "application/mpeg", blob)
}

// POST /api/v1/upload
func (r *router) uploadFile(c echo.Context) error {
	name := c.FormValue("name")
	multipart, err := c.FormFile("file")
	if err != nil {
		return r.errHandler.handle(err)
	}

	return r.upload(c, name, multipart)
}

func (r *router) upload(c echo.Context, name string, multipart *multipart.FileHeader) error {
	ctx := c.Request().Context()

	// check if attached file is of mp4 format
	if !strings.HasSuffix(multipart.Filename, ".mp4") {
		return echo.ErrUnprocessableEntity
	}

	video, err := multipart.Open()
	if err != nil {
		return r.errHandler.handle(err)
	}

	// uploading all the created files
	if err := r.service.Upload(ctx, video, name); err != nil {
		return r.errHandler.handle(err)
	}

	return c.JSON(200, "Success")
}

// DELETE /api/v1/:filename
func (r *router) deleteFile(c echo.Context) error {
	filename := c.Param("filename")

	return r.delete(c, filename)
}

func (r *router) delete(c echo.Context, filename string) error {
	ctx := c.Request().Context()

	if err := r.service.Remove(ctx, filename); err != nil {
		return r.errHandler.handle(err)
	}

	return c.JSON(200, "Success")
}
