package controller

import (
	"io"
	"log"
	"mime/multipart"
	"strings"

	"github.com/cutlery47/gostream/internal/schema"
	"github.com/cutlery47/gostream/internal/service"
	"github.com/cutlery47/gostream/internal/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type router struct {
	service    service.FileService
	errHandler errHandler
}

func newRouter(errLog *zap.Logger, service service.FileService) *router {
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
	var file *schema.OutFile

	// searching for requested file
	file, err = r.service.Serve(filename)
	if err != nil {
		return r.errHandler.handle(err)
	}

	// converting the file into a sequence of bytes
	blob, err := utils.BufferReader(file.Raw)
	if err != nil {
		return r.errHandler.handle(err)
	}

	// returning the file
	return c.Blob(200, "application/mpeg", blob.Bytes())
}

// POST /api/v1/upload
func (r *router) uploadFile(c echo.Context) error {
	filename := c.FormValue("filename")
	multipart, err := c.FormFile("file")
	if err != nil {
		return r.errHandler.handle(err)
	}

	return r.upload(c, filename, multipart)
}

func (r *router) upload(c echo.Context, filename string, multipart *multipart.FileHeader) error {
	// check if attached file is of mp4 format
	if !strings.HasSuffix(multipart.Filename, ".mp4") {
		return echo.ErrUnprocessableEntity
	}

	video, err := multipart.Open()
	if err != nil {
		return r.errHandler.handle(err)
	}

	inVideo := schema.InFile{
		Raw:  video,
		Name: filename,
		Size: int(multipart.Size),
	}

	// uploading all the created files
	if err := r.service.Upload(inVideo); err != nil {
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
	if err := r.service.Remove(filename); err != nil {
		return r.errHandler.handle(err)
	}

	return c.JSON(200, "Success")
}

func (r *router) getMinioLogs(c echo.Context) error {
	body := c.Request().Body

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return r.errHandler.handle(err)
	}

	log.Println(string(bodyBytes))

	return c.JSON(200, "Success")
}
