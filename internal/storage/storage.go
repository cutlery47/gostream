package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cutlery47/gostream/internal/utils"
	"go.uber.org/zap"
)

// abstracts out file manipulation
type Storage interface {
	// stores files
	Store(ctx context.Context, video, manifest File, chunks []File) error
	// retrieves file
	Get(ctx context.Context, filename string) (io.ReadCloser, error)
	// removes file
	Remove(ctx context.Context, filename string) error
	// returns local paths
	Paths() Paths
}

// db + obj storage based storage
type DistibutedStorage struct {
	repo Repository
	s3   ObjectStorage

	infoLog *zap.Logger
	errLog  *zap.Logger
	paths   Paths
}

func NewDistibutedStorage(infoLog, errLog *zap.Logger, paths Paths, repo Repository, s3 ObjectStorage) *DistibutedStorage {
	return &DistibutedStorage{
		repo:    repo,
		s3:      s3,
		infoLog: infoLog,
		errLog:  errLog,
		paths:   paths,
	}
}

// todo: make s3 uploads "transactional"
func (ds *DistibutedStorage) Store(ctx context.Context, video, manifest File, chunks []File) error {
	vidLocation, err := ds.s3.Store(ctx, video)
	if err != nil {
		return err
	}

	manLocation, err := ds.s3.Store(ctx, manifest)
	if err != nil {
		return err
	}

	chunkLocations, err := ds.s3.StoreMultiple(ctx, chunks...)
	if err != nil {
		return err
	}

	video.Location = vidLocation
	manifest.Location = manLocation
	for i := range chunks {
		chunks[i].Location = chunkLocations[i]
	}

	return ds.repo.CreateAll(ctx, video, manifest, chunks)
}

func (ds *DistibutedStorage) Get(ctx context.Context, filename string) (io.ReadCloser, error) {
	fileLocation, err := ds.repo.Read(ctx, filename)
	if err != nil {
		return nil, err
	}

	return ds.s3.Get(ctx, fileLocation)
}

func (ds *DistibutedStorage) Remove(ctx context.Context, filename string) error {
	fileLocation, err := ds.repo.Delete(ctx, filename)
	if err != nil {
		return err
	}

	return ds.s3.Delete(ctx, fileLocation)
}

func (ds DistibutedStorage) Paths() Paths {
	return ds.paths
}

// local file system based storage
type LocalStorage struct {
	errLog *zap.Logger
	paths  Paths
}

func NewLocalStorage(errLog *zap.Logger, paths Paths) *LocalStorage {
	return &LocalStorage{
		errLog: errLog,
		paths:  paths,
	}
}

func (ls *LocalStorage) Store(ctx context.Context, video, manifest File, chunks []File) error {
	// when storing files locally, there is no need to write file to any other storage
	return nil
}

func (ls *LocalStorage) Get(ctx context.Context, filename string) (io.ReadCloser, error) {
	filePath, err := ls.determinePath(filename)
	if err != nil {
		return nil, err
	}

	return os.Open(filePath)
}

func (ls *LocalStorage) Remove(ctx context.Context, filename string) error {
	filePath, err := ls.determinePath(filename)
	if err != nil {
		return err
	}

	return os.Remove(filePath)
}

func (ls *LocalStorage) Paths() Paths {
	return ls.paths
}

// used to detect where given file is stored
func (ls *LocalStorage) determinePath(filename string) (filePath string, err error) {
	if strings.HasSuffix(filename, ".mp4") {
		filePath = fmt.Sprintf("%v/%v", ls.paths.VidPath, filename)
	} else if strings.HasSuffix(filename, ".m3u8") {
		filePath = fmt.Sprintf("%v/%v", ls.paths.ManPath, filename)
	} else if strings.HasSuffix(filename, ".ts") {
		subdir := utils.RemoveSuffix(filename, "_")
		filePath = fmt.Sprintf("%v/%v/%v", ls.paths.ChunkPath, subdir, filename)
	} else {
		return filePath, ErrUnsupportedFileFormat
	}

	return filePath, nil
}

// structure for storing paths to local files
type Paths struct {
	VidPath   string
	ManPath   string
	ChunkPath string
}
