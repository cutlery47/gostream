package obj

import "io"

type S3File struct {
	file *io.Reader
}
