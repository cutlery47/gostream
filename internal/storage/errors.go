package storage

import "errors"

var (
	ErrNotImplemented        = errors.New("feature is not yet implemented")
	ErrUnsupportedFileFormat = errors.New("unsupported file format")
	ErrUniueVideo            = errors.New("video with provided name already exists")
)
