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
	Remove(filename string) error
	// returns a path to the storage
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

func (lms *LocalManifestStorage) Get(filename string) (io.Reader, error) {
	return os.Open(fmt.Sprintf("%v/%v.m3u8", lms.manifestPath, filename))
}

func (lms *LocalManifestStorage) Exists(filename string) bool {
	if _, err := os.Stat(fmt.Sprintf("%v/%v.m3u8", lms.manifestPath, filename)); err == nil {
		return true
	}
	return false
}

func (lms *LocalManifestStorage) Store(file io.Reader, filename string) error {
	return ErrNotImplemented
}

func (lms *LocalManifestStorage) Remove(filename string) error {
	return os.Remove(fmt.Sprintf("%v/%v.m3u8", lms.manifestPath, filename))
}

func (lms *LocalManifestStorage) Path() string {
	return lms.manifestPath
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

func (lcs *LocalChunkStorage) Exists(filename string) bool {
	chunkdir := utils.RemoveSuffix(filename, "_")
	if _, err := os.Stat(fmt.Sprintf("%v/%v/%v", lcs.chunkPath, chunkdir, filename)); err == nil {
		return true
	}
	return false
}

func (lcs *LocalChunkStorage) Store(file io.Reader, filename string) error {
	return ErrNotImplemented
}

func (lcs *LocalChunkStorage) Remove(filename string) error {
	chunkdir := utils.RemoveSuffix(filename, "_")
	return os.Remove(fmt.Sprintf("%v/%v/%v", lcs.chunkPath, chunkdir, filename))
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

func (lvs *LocalVideoStorage) Exists(filename string) bool {
	if _, err := os.Stat(fmt.Sprintf("%v/%v", lvs.videoPath, filename)); err == nil {
		return true
	}
	return false
}

func (lvs *LocalVideoStorage) Store(file io.Reader, filename string) error {
	rawFile, err := utils.BufferReader(file)
	if err != nil {
		return err
	}

	newFile, err := os.Create(fmt.Sprintf("%v/%v", lvs.videoPath, filename))
	if err != nil {
		return err
	}

	_, err = io.Copy(newFile, rawFile)
	if err != nil {
		return err
	}

	return nil
}

func (lvs *LocalVideoStorage) Remove(filename string) error {
	return os.Remove(fmt.Sprintf("%v/%v", lvs.videoPath, filename))
}

func (lvs *LocalVideoStorage) Path() string {
	return lvs.videoPath
}
