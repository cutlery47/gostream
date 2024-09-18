package service

import "os"

type Service interface {
	Serve(filename string) (*os.File, error)
}

type chunkService struct {
}

func NewChunkService() *chunkService {
	return &chunkService{}
}

func (cs *chunkService) Serve(filename string) (*os.File, error) {
	return nil, nil
}

type manifestService struct {
}

func NewManifestService() *manifestService {
	return &manifestService{}
}

func (ms *manifestService) Serve(filename string) (*os.File, error) {
	return nil, nil
}
