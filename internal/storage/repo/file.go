package repo

import "time"

type InRepositoryFile struct {
	Name       string
	Size       int
	UploadedAt time.Time
	BucketName string
	ETag       string
}

type RepositoryFile struct {
	id int
	InRepositoryFile
}
