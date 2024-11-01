package storage

import (
	"io"
	"os"
)

type File struct {
	// binary file reader
	Raw io.ReadCloser
	// name for database
	FileName string
	// name for object storage
	ObjectName string
	// file location in the obj storage
	Location Location

	Size int64
}

func FromFD(file *os.File, filename string) (*File, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &File{
		Raw:        file,
		FileName:   filename,
		ObjectName: file.Name(),
		Location:   Location{},
		Size:       info.Size(),
	}, nil
}

type Location struct {
	// object storage bucket
	Bucket string
	// object storage key
	Object string
}
