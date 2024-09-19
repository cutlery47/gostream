package storage

import (
	"fmt"
	"os"
)

type Storage interface {
	Get(filename string) (*os.File, error)
	Exists(filename string) bool
	Path() string
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
	return os.Open(fmt.Sprintf("%v/%v", lfs.manifestPath, filename))
}

func (lfs *LocalManifestStorage) Exists(filename string) bool {
	if _, err := os.Stat(fmt.Sprintf("%v/%v", lfs.manifestPath, filename)); err == nil {
		return true
	}
	return false
}

func (lfs *LocalManifestStorage) Path() string {
	return lfs.manifestPath
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

func (lcs *LocalChunkStorage) Exists(filename string) bool {
	if _, err := os.Stat(fmt.Sprintf("%v/%v", lcs.chunkPath, filename)); err == nil {
		return true
	}
	return false
}

func (lcs *LocalChunkStorage) Path() string {
	return lcs.chunkPath
}

type LocalVideoStorage struct {
	videoPath string
}

func NewLocalVideoStorage(videoPath string) *LocalVideoStorage {
	return &LocalVideoStorage{
		videoPath: videoPath,
	}
}

func (lvs *LocalVideoStorage) Get(filename string) (*os.File, error) {
	return os.Open(fmt.Sprintf("%v/%v", lvs.videoPath, filename))
}

func (lvs *LocalVideoStorage) Exists(filename string) bool {
	if _, err := os.Stat(fmt.Sprintf("%v/%v", lvs.videoPath, filename)); err == nil {
		return true
	}
	return false
}

func (lvs *LocalVideoStorage) Path() string {
	return lvs.videoPath
}
