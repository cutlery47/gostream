package v1

import (
	"io"
	"strings"

	"github.com/cutlery47/gostream/internal/service"
	"github.com/labstack/echo/v4"
)

type fileRoutes struct {
	s service.Service
	h *errHandler
}

func newFileRoutes(g *echo.Group, s service.Service, h *errHandler) {
	r := &fileRoutes{
		s: s,
		h: h,
	}

	g.POST("/", r.upload)
	g.GET("/:filename", r.get)
	g.DELETE("/:filename", r.delete)
}

// @Summary Deposit
// @Description Deposit
// @Tags accounts
// @Accept json
// @Produce json
// @Param input body v1.accountDepositInput true "input"
// @Success 200
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Security JWT
// @Router /api/v1/accounts/deposit [post]
func (r *fileRoutes) get(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	filename := c.Param("filename")

	var file io.ReadCloser

	ctx := c.Request().Context()

	// searching for requested file
	file, err := r.s.Serve(ctx, filename)
	if err != nil {
		return r.h.handle(err)
	}

	// converting the file into a sequence of bytes
	blob, err := io.ReadAll(file)
	if err != nil {
		return r.h.handle(err)
	}

	// returning the file
	return c.Blob(200, "application/mpeg", blob)
}

// POST /api/v1/upload
func (r *fileRoutes) upload(c echo.Context) error {
	name := c.FormValue("name")
	multipart, err := c.FormFile("file")
	if err != nil {
		return r.h.handle(err)
	}

	ctx := c.Request().Context()

	// check if attached file is of mp4 format
	if !strings.HasSuffix(multipart.Filename, ".mp4") {
		return echo.ErrUnprocessableEntity
	}

	video, err := multipart.Open()
	if err != nil {
		return r.h.handle(err)
	}

	// uploading all the created files
	if err := r.s.Upload(ctx, video, name); err != nil {
		return r.h.handle(err)
	}

	return c.JSON(200, "Success")
}

// DELETE /api/v1/:filename
func (r *fileRoutes) delete(c echo.Context) error {
	filename := c.Param("filename")

	ctx := c.Request().Context()

	if err := r.s.Remove(ctx, filename); err != nil {
		return r.h.handle(err)
	}

	return c.JSON(200, "Success")
}
