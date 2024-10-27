package obj

import (
	"io"

	"github.com/cutlery47/gostream/internal/schema"
)

type InFile struct {
	Raw  io.ReadCloser
	Name string
	Size int
}

func FromSchema(file schema.InFile) InFile {
	return InFile{
		Raw:  file.Raw,
		Name: file.Name,
		Size: file.Size,
	}
}

type S3FileInfo struct {
}
