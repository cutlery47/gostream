package controller

import (
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/cutlery47/gostream/internal/service"
	"github.com/cutlery47/gostream/internal/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type router struct {
	videoService    service.UploadRemoveService
	manifestService service.Service
	chunkService    service.Service
	errHandler      errHandler
}

func newRouter(
	errLog *zap.Logger,
	chunkService service.Service,
	manifestService service.Service,
	videoService service.UploadRemoveService) *router {

	return &router{
		manifestService: manifestService,
		chunkService:    chunkService,
		videoService:    videoService,
		errHandler:      *newErrHandler(errLog),
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
	file, err := c.FormFile("file")
	if err != nil {
		return r.errHandler.handle(err)
	}

	multipart, err := file.Open()
	if err != nil {
		return r.errHandler.handle(err)
	}
	defer multipart.Close()

	return r.upload(c, file.Filename, multipart)
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
		file, err = r.chunkService.Serve(filename)
		if err != nil {
			return r.errHandler.handle(err)
		}
	} else if strings.HasSuffix(filename, ".mp4") {
		// entire file was requested
		file, err = r.videoService.Serve(filename)
		if err != nil {
			return r.errHandler.handle(err)
		}
	} else {
		// manifest file was requested
		file, err = r.manifestService.Serve(filename)
		if err != nil {
			return r.errHandler.handle(err)
		}
	}

	// converting the file into a sequence of bytes
	blob, err := utils.BufferReader(file)
	if err != nil {
		return r.errHandler.handle(err)
	}

	// returning the file
	return c.Blob(200, "application/mpeg", blob.Bytes())
}

func (r *router) upload(c echo.Context, filename string, multipart multipart.File) error {
	if strings.HasSuffix(filename, ".mp4") {
		if err := r.videoService.Upload(multipart, filename); err != nil {
			return r.errHandler.handle(err)
		}
	} else {
		return echo.ErrUnprocessableEntity
	}

	return c.JSON(200, "Success")
}

func (r *router) delete(c echo.Context, filename string) error {
	if strings.HasSuffix(filename, ".mp4") {
		if err := r.videoService.Remove(filename); err != nil {
			return r.errHandler.handle(err)
		}
	} else {
		if err := r.videoService.Remove(fmt.Sprintf("%v.mp4", filename)); err != nil {
			return r.errHandler.handle(err)
		}
	}

	return c.JSON(200, "Success")
}
