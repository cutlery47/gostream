package obj

import (
	"fmt"
	"net"

	"github.com/cutlery47/gostream/config"
)

type ObjectStorage interface {
	Store(file S3File) error
	Get(filename string) (*S3File, error)
}

type S3 struct {
	conn net.Conn
}

func (s3 S3) Store(file S3File) error {
	return fmt.Errorf("fdsfd")
}

func (s3 S3) Get(filename string) (*S3File, error) {
	return nil, fmt.Errorf("sdafasdf")
}

func NewS3(conf config.S3Config) (*S3, error) {
	return &S3{}, nil
}
