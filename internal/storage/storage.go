package storage

import (
	"fmt"
	"os"
)

type Storage interface {
	Get(filename string) (*os.File, error)
}

type LocalManifestStorage struct {
	manifestPath string
}

func NewLocalManifestStorage(manifestPath string) *LocalManifestStorage {
	return &LocalManifestStorage{
		manifestPath: manifestPath,
	}
}

func (lfs *LocalManifestStorage) Get(filename string) (*os.File, error) {
	return os.Open(fmt.Sprintf("%v/%v.m3u8", lfs.manifestPath, filename))
}

type LocalChunkStorage struct {
	chunkPath string
}

func NewLocalChunkStorage(chunkPath string) *LocalChunkStorage {
	return &LocalChunkStorage{
		chunkPath: chunkPath,
	}
}

func (lcs *LocalChunkStorage) Get(filename string) (*os.File, error) {
	return nil, nil
}
