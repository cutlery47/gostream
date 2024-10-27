package storage

import (
	"fmt"
	"os"
	"strings"

	"github.com/cutlery47/gostream/internal/schema"
	"github.com/cutlery47/gostream/internal/storage/obj"
	"github.com/cutlery47/gostream/internal/storage/repo"
	"github.com/cutlery47/gostream/internal/utils"
	"go.uber.org/zap"
)

// abstracts out file manipulation
type Storage interface {
	// stores files
	Store(video schema.InVideo, manifest schema.InFile, chunks []schema.InFile) error
	// retrieves file
	Get(filename string) (*schema.OutFile, error)
	// removes file
	Remove(filename string) error
	// checks if file exists
	Exists(filename string) bool
	// returns local paths
	Paths() Paths
}

// db + obj storage based storage
type DistibutedStorage struct {
	repo repo.Repository
	s3   obj.ObjectStorage

	infoLog *zap.Logger
	errLog  *zap.Logger
	paths   Paths
}

func NewDistibutedStorage(infoLog, errLog *zap.Logger, paths Paths, repo repo.Repository, s3 obj.ObjectStorage) *DistibutedStorage {
	return &DistibutedStorage{
		repo:    repo,
		s3:      s3,
		infoLog: infoLog,
		errLog:  errLog,
		paths:   paths,
	}
}

// todo: make s3 uploads "transactional"
func (ds *DistibutedStorage) Store(video schema.InVideo, manifest schema.InFile, chunks []schema.InFile) error {
	fmt.Printf("vid: %+v\nman: %+v\nchunks: %+v\n", video, manifest, chunks)

	vidS3 := obj.FromSchema(video.File)
	vidLocation, err := ds.s3.Store(vidS3)
	if err != nil {
		return err
	}

	manS3 := obj.FromSchema(manifest)
	manLocation, err := ds.s3.Store(manS3)
	if err != nil {
		return err
	}

	chunksS3 := []obj.InFile{}
	for _, chunk := range chunks {
		chunksS3 = append(chunksS3, obj.FromSchema(chunk))
	}
	chunksLocations, err := ds.s3.StoreMultiple(chunksS3...)
	if err != nil {
		return err
	}

	repoVid := video.ToRepo(vidLocation)
	repoMan := manifest.ToRepo(manLocation)
	repoChunks := []repo.InFile{}
	for i, el := range chunks {
		repoChunks = append(repoChunks, el.ToRepo(chunksLocations[i]))
	}

	// remove local files here

	return ds.repo.CreateAll(repoVid, repoMan, repoChunks)
}

func (ds *DistibutedStorage) Get(filename string) (*schema.OutFile, error) {
	repoFile, err := ds.repo.Read(filename)
	if err != nil {
		return nil, err
	}

	s3File, err := ds.s3.Get(repoFile.Data.Name)
	if err != nil {
		return nil, err
	}

	return &schema.OutFile{Raw: s3File.Raw}, nil
}

func (ds *DistibutedStorage) Remove(filename string) error {
	repoFile, err := ds.repo.Delete(filename)
	if err != nil {
		return err
	}

	if _, err := ds.s3.Delete(repoFile.Data.Name); err != nil {
		return err
	}

	return nil
}

func (ds *DistibutedStorage) Exists(filename string) bool {
	repoFile, err := ds.repo.Read(filename)
	if err != nil {
		return false
	}

	if _, err := ds.s3.Get(repoFile.Data.Name); err != nil {
		return false
	}

	return true
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

func (ls *LocalStorage) Store(video schema.InVideo, manifest schema.InFile, chunks []schema.InFile) error {
	// when storing files locally, there is no need to write file to any other storage
	return nil
}

func (ls *LocalStorage) Get(filename string) (*schema.OutFile, error) {
	filePath, err := ls.determinePath(filename)
	if err != nil {
		return nil, err
	}

	fd, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &schema.OutFile{Raw: fd}, nil
}

func (ls *LocalStorage) Remove(filename string) error {
	filePath, err := ls.determinePath(filename)
	if err != nil {
		return err
	}

	return os.Remove(filePath)
}

func (ls *LocalStorage) Exists(filename string) bool {
	filePath, err := ls.determinePath(filename)
	if err != nil {
		ls.errLog.Error(err.Error())
		return false
	}

	if _, err := os.Stat(filePath); err == nil {
		return true
	}

	return false
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
