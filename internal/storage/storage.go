package storage

import (
	"fmt"
	"io"
	"os"

	"github.com/cutlery47/gostream/internal/utils"
)

type Storage interface {
	Get(filename string) (io.Reader, error)
	Exists(filename string) bool
	Store(file io.Reader, filename string) error
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

func (lfs *LocalManifestStorage) Get(filename string) (io.Reader, error) {
	return os.Open(fmt.Sprintf("%v/%v", lfs.manifestPath, filename))
}

func (lfs *LocalManifestStorage) Store(file io.Reader, filename string) error {
	return ErrNotImplemented
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

func (lcs *LocalChunkStorage) Get(filename string) (io.Reader, error) {
	chunkdir := utils.RemoveSuffix(filename, "_")
	return os.Open(fmt.Sprintf("%v/%v/%v", lcs.chunkPath, chunkdir, filename))
}

func (lcs *LocalChunkStorage) Store(file io.Reader, filename string) error {
	return ErrNotImplemented
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

func (lvs *LocalVideoStorage) Get(filename string) (io.Reader, error) {
	return os.Open(fmt.Sprintf("%v/%v", lvs.videoPath, filename))
}

func (lcs *LocalVideoStorage) Store(file io.Reader, filename string) error {
	rawFile, err := utils.BufferReader(file)
	if err != nil {
		return err
	}

	newFile, err := os.Create(fmt.Sprintf("%v/%v", lcs.videoPath, filename))
	if err != nil {
		return err
	}

	_, err = io.Copy(newFile, rawFile)
	if err != nil {
		return err
	}

	return nil
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
