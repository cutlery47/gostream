package repo

import "errors"

var (
	ErrUniqueVideo = errors.New("video with provided name already exists")
)
