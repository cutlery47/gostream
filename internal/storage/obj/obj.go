package obj

import (
	"fmt"
	"net"

	"github.com/cutlery47/gostream/config"
)

type ObjectStorage interface {
	Store(file S3File) (string, string, error)
	StoreMultiple(files ...S3File) ([]string, []string, error)
	Get(filename string) (*S3File, error)
	Delete(filename string) (*S3File, error)
}

type S3 struct {
	conn net.Conn
}

func NewS3(conf config.S3Config) (*S3, error) {
	return &S3{}, nil
}

func (s3 S3) Store(files S3File) (string, string, error) {
	return "", "", fmt.Errorf("fdsfd")
}

func (s3 S3) Get(filename string) (*S3File, error) {
	return nil, fmt.Errorf("sdafasdf")
}

func (s3 S3) Delete(filename string) (*S3File, error) {
	return nil, fmt.Errorf("xyu yxu")
}

func (s3 S3) StoreMultiple(file ...S3File) ([]string, []string, error) {
	return []string{}, []string{}, fmt.Errorf("sfkasdf")
}
