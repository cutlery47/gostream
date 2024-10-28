package schema

import (
	"io"

	"github.com/cutlery47/gostream/internal/storage/repo"
)

type InFile struct {
	// binary file reader
	Raw io.ReadCloser
	// initial file name
	Name string
	// file size
	Size int
}

type InVideo struct {
	File InFile
	// video name, provided by user
	Name string
}

func (f InFile) ToRepo(location string) repo.InFile {
	return repo.InFile{
		Name:     f.Name,
		Size:     f.Size,
		Location: location,
	}
}

func (v InVideo) ToRepo(location string) repo.InVideo {
	return repo.InVideo{
		File:      v.File.ToRepo(location),
		VideoName: v.Name,
	}
}

type OutFile struct {
	Raw io.ReadCloser
}
