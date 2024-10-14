package storage

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cutlery47/gostream/internal/schema"
	"github.com/cutlery47/gostream/internal/storage/obj"
	"github.com/cutlery47/gostream/internal/storage/repo"
	"github.com/cutlery47/gostream/internal/utils"
)

type Storage interface {
	// returns file by name
	Get(fileName string) (*schema.OutFile, error)
	// checks if file exists
	Exists(fileName string) bool
	// stores file
	Store(file schema.InFile) error
	// removes file by name
	Remove(fileName string) error
	// returns a path to the storage
	Path() string
}

type DistributedManifestStorage struct {
	repo repo.CreateReadDeleteRepository
	s3   obj.ObjectStorage
}

func NewDistributedManifestStorage(repository repo.CreateReadDeleteRepository, s3 obj.ObjectStorage) *DistributedManifestStorage {
	return &DistributedManifestStorage{
		repo: repository,
		s3:   s3,
	}
}

func (dms DistributedManifestStorage) Store(file schema.InFile) error {
	repoFile := repo.InRepositoryFile{
		Name:       file.Name,
		Size:       file.Size,
		UploadedAt: time.Now(),
		BucketName: "tmp",
		ETag:       "xyu",
	}

	if err := dms.repo.Create(repoFile); err != nil {
		return err
	}

	return nil
}

func (dms DistributedManifestStorage) Get(fileName string) (*schema.OutFile, error) {
	return nil, ErrNotImplemented
}

func (dms DistributedManifestStorage) Exists(fileName string) bool {
	return false
}

func (dms DistributedManifestStorage) Remove(fileName string) error {
	return ErrNotImplemented
}

func (dms DistributedManifestStorage) Path() string {
	return ""
}

type DistributedVideoStorage struct {
	repo repo.CreateReadDeleteRepository
	s3   obj.ObjectStorage
}

func NewDistributedVideoStorage(repository repo.CreateReadDeleteRepository, s3 obj.ObjectStorage) *DistributedVideoStorage {
	return &DistributedVideoStorage{
		repo: repository,
		s3:   s3,
	}
}

func (dvs DistributedVideoStorage) Store(file schema.InFile) error {
	repoFile := repo.InRepositoryFile{
		Name:       file.Name,
		Size:       file.Size,
		UploadedAt: time.Now(),
		BucketName: "tmp",
		ETag:       "xyu",
	}

	if err := dvs.repo.Create(repoFile); err != nil {
		return err
	}

	return nil
}

func (dvs DistributedVideoStorage) Get(fileName string) (*schema.OutFile, error) {
	return nil, ErrNotImplemented
}

func (dvs DistributedVideoStorage) Exists(fileName string) bool {
	return false
}

func (dvs DistributedVideoStorage) Remove(fileName string) error {
	return ErrNotImplemented
}

func (dvs DistributedVideoStorage) Path() string {
	return ""
}

type LocalManifestStorage struct {
	manifestPath string
}

type DistributedChunkStorage struct {
	repo repo.CreateReadDeleteRepository
	s3   obj.ObjectStorage
}

func NewDistributedChunkStorage(repository repo.CreateReadDeleteRepository, s3 obj.ObjectStorage) *DistributedChunkStorage {
	return &DistributedChunkStorage{
		repo: repository,
		s3:   s3,
	}
}

func (dcs DistributedChunkStorage) Store(file schema.InFile) error {
	repoFile := repo.InRepositoryFile{
		Name:       file.Name,
		Size:       file.Size,
		UploadedAt: time.Now(),
		BucketName: "tmp",
		ETag:       "xyu",
	}

	if err := dcs.repo.Create(repoFile); err != nil {
		return err
	}

	return nil
}

func (dvs DistributedChunkStorage) Get(fileName string) (*schema.OutFile, error) {
	return nil, ErrNotImplemented
}

func (dvs DistributedChunkStorage) Exists(fileName string) bool {
	return false
}

func (dvs DistributedChunkStorage) Remove(fileName string) error {
	return ErrNotImplemented
}

func (dvs DistributedChunkStorage) Path() string {
	return ""
}

func NewLocalManifestStorage(manifestPath string) *LocalManifestStorage {
	return &LocalManifestStorage{
		manifestPath: manifestPath,
	}
}

func (lms *LocalManifestStorage) Get(fileName string) (*schema.OutFile, error) {
	fd, err := os.Open(fmt.Sprintf("%v/%v", lms.manifestPath, fileName))
	if err != nil {
		return nil, err
	}

	return &schema.OutFile{Raw: fd}, nil
}

func (lms *LocalManifestStorage) Exists(filename string) bool {
	if _, err := os.Stat(fmt.Sprintf("%v/%v", lms.manifestPath, filename)); err == nil {
		return true
	}
	return false
}

func (lms *LocalManifestStorage) Store(file schema.InFile) error {
	return ErrNotImplemented
}

func (lms *LocalManifestStorage) Remove(fileName string) error {
	return os.Remove(fmt.Sprintf("%v/%v", lms.manifestPath, fileName))
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

func (lcs *LocalChunkStorage) Get(filename string) (*schema.OutFile, error) {
	chunkdir := utils.RemoveSuffix(filename, "_")
	fd, err := os.Open(fmt.Sprintf("%v/%v/%v", lcs.chunkPath, chunkdir, filename))
	if err != nil {
		return nil, err
	}

	return &schema.OutFile{Raw: fd}, nil
}

func (lcs *LocalChunkStorage) Exists(filename string) bool {
	chunkdir := utils.RemoveSuffix(filename, "_")
	if _, err := os.Stat(fmt.Sprintf("%v/%v/%v", lcs.chunkPath, chunkdir, filename)); err == nil {
		return true
	}
	return false
}

func (lcs *LocalChunkStorage) Store(file schema.InFile) error {
	return ErrNotImplemented
}

func (lcs *LocalChunkStorage) Remove(filename string) error {
	return os.RemoveAll(fmt.Sprintf("%v/%v", lcs.chunkPath, filename))
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

func (lvs *LocalVideoStorage) Get(filename string) (*schema.OutFile, error) {
	fd, err := os.Open(fmt.Sprintf("%v/%v", lvs.videoPath, filename))
	if err != nil {
		return nil, err
	}

	return &schema.OutFile{Raw: fd}, nil
}

func (lvs *LocalVideoStorage) Exists(filename string) bool {
	if _, err := os.Stat(fmt.Sprintf("%v/%v", lvs.videoPath, filename)); err == nil {
		return true
	}
	return false
}

func (lvs *LocalVideoStorage) Store(file schema.InFile) error {
	rawFile, err := utils.BufferReader(file.Raw)
	if err != nil {
		return err
	}

	// creating a new file on a local system
	newFile, err := os.Create(fmt.Sprintf("%v/%v", lvs.videoPath, file.Name))
	if err != nil {
		return err
	}

	// copying data to the file
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
