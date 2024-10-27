package repo

import "time"

type InFile struct {
	Name       string
	Size       int
	UploadedAt time.Time
	// location of the file in s3
	Location string
}

type InVideo struct {
	File      InFile
	VideoName string
}

type File struct {
	Id   int
	Data InFile
}
