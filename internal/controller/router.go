package controller

import (
	"io"
	"mime/multipart"
	"strings"

	"github.com/cutlery47/gostream/internal/service"
	"github.com/cutlery47/gostream/internal/utils"
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

// POST /api/v1/upload
func (r *router) uploadFile(c echo.Context) error {
	multipart, err := c.FormFile("file")
	if err != nil {
		return r.errHandler.handle(err)
	}

	filename := c.FormValue("filename")

	return r.upload(c, filename, multipart)
}

// DELETE /api/v1/:filename
func (r *router) deleteFile(c echo.Context) error {
	filename := c.Param("filename")

	return r.delete(c, filename)
}

func (r *router) get(c echo.Context, filename string) (err error) {
	var file io.Reader

	// searching for requested file on the current system
	if strings.HasSuffix(filename, ".ts") {
		// transport stream was requested
		file, err = r.service.ServeChunk(filename)
	} else if strings.HasSuffix(filename, ".mp4") {
		// entire file was requested
		file, err = r.service.ServeEntireVideo(filename)
	} else {
		// manifest file was requested
		file, err = r.service.ServeManifest(filename)
	}

	if err != nil {
		return r.errHandler.handle(err)
	}

	// converting the file into a sequence of bytes
	blob, err := utils.BufferReader(file)
	if err != nil {
		return r.errHandler.handle(err)
	}

	// returning the file
	return c.Blob(200, "application/mpeg", blob.Bytes())
}

func (r *router) upload(c echo.Context, filename string, multipart *multipart.FileHeader) error {
	// check if attached file is of mp4 format
	if !strings.HasSuffix(multipart.Filename, ".mp4") {
		return echo.ErrUnprocessableEntity
	}

	file, err := multipart.Open()
	if err != nil {
		return err
	}

	filename += ".mp4"

	if err := r.service.UploadVideo(file, filename); err != nil {
		return r.errHandler.handle(err)
	}

	return c.JSON(200, "Success")
}

func (r *router) delete(c echo.Context, filename string) error {
	if err := r.service.RemoveVideo(filename); err != nil {
		return r.errHandler.handle(err)
	}

	return c.JSON(200, "Success")
}
