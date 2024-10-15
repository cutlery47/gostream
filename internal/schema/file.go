package schema

import (
	"io"
	"time"

	"github.com/cutlery47/gostream/internal/storage/repo"
)

type InFile struct {
	Raw  io.Reader
	Name string
	Size int
}

func (f InFile) ToRepo(bucketName, eTag string) repo.InRepositoryFile {
	return repo.InRepositoryFile{
		Name:       f.Name,
		Size:       f.Size,
		UploadedAt: time.Now(),
		BucketName: bucketName,
		ETag:       eTag,
	}
}

type OutFile struct {
	Raw io.Reader
}
