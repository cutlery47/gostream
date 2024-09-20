package controller

import (
	"strings"

	"github.com/cutlery47/gostream/internal/service"
	"github.com/cutlery47/gostream/internal/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type router struct {
	manifestService service.Service
	chunkService    service.Service
	videoService    service.Service
	errHandler      errHandler
}

func newRouter(manifestService, chunkService, videoService service.Service, errLog *zap.Logger) *router {
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
	if strings.HasSuffix(filename, ".ts") {
		// transport stream was requested
		return r.serveFile(c, filename, r.chunkService)
	} else {
		// manifest file was requested
		return r.serveFile(c, filename, r.manifestService)
	}
}

func (r *router) serveFile(c echo.Context, filename string, service service.Service) error {
	// searching for requested file on the current system
	file, err := service.Serve(filename)
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

	if strings.HasSuffix(file.Filename, ".mp4") {
		r.videoService.Upload(multipart)
	}

	return echo.ErrNotImplemented
}
