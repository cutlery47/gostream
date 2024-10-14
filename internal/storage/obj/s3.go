package obj

import "io"

type S3 interface {
	Store(file io.Reader)
	Get(filename string) S3File
}
