package storage

import (
	"fmt"
	"os"
	"strings"

	"github.com/cutlery47/gostream/internal/schema"
	"github.com/cutlery47/gostream/internal/storage/obj"
	"github.com/cutlery47/gostream/internal/storage/repo"
	"go.uber.org/zap"
)

type Storage interface {
	// stores files
	Store(video, manifest schema.InFile, chunks []schema.InFile) error
	// retrieves file
	Get(filename string) (*schema.OutFile, error)
	// removes file
	Remove(filename string) error
	// checks if file exists
	Exists(filename string) bool
}

type DistibutedStorage struct {
	log  *zap.Logger
	repo repo.Repository
	s3   obj.ObjectStorage
}

func NewDistibutedStorage(log *zap.Logger, repo repo.Repository, s3 obj.ObjectStorage) *DistibutedStorage {
	return &DistibutedStorage{
		log:  log,
		repo: repo,
		s3:   s3,
	}
}

// make s3 uploads "transaactional"
func (ds *DistibutedStorage) Store(video, manifest schema.InFile, chunks []schema.InFile) error {
	vidS3 := obj.S3File{Raw: video.Raw}
	vidBucketName, vidETag, err := ds.s3.Store(vidS3)
	if err != nil {
		return err
	}

	manS3 := obj.S3File{Raw: video.Raw}
	manBucketName, manETag, err := ds.s3.Store(manS3)
	if err != nil {
		return err
	}

	chunksS3 := []obj.S3File{}
	for _, el := range chunks {
		chunksS3 = append(chunksS3, obj.S3File{Raw: el.Raw})
	}
	chunksBucketNames, chunksETags, err := ds.s3.StoreMultiple(chunksS3...)
	if err != nil {
		return err
	}

	repoVid := video.ToRepo(vidBucketName, vidETag)
	repoMan := manifest.ToRepo(manBucketName, manETag)
	repoChunks := []repo.InRepositoryFile{}
	for i, el := range chunks {
		repoChunks = append(repoChunks, el.ToRepo(chunksBucketNames[i], chunksETags[i]))
	}

	return ds.repo.CreateAll(repoVid, repoMan, repoChunks)
}

func (ds *DistibutedStorage) Get(filename string) (*schema.OutFile, error) {
	repoFile, err := ds.repo.Read(filename)
	if err != nil {
		return nil, err
	}

	s3File, err := ds.s3.Get(repoFile.Name)
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

	if _, err := ds.s3.Delete(repoFile.Name); err != nil {
		return err
	}

	return nil
}

func (ds *DistibutedStorage) Exists(filename string) bool {
	repoFile, err := ds.repo.Read(filename)
	if err != nil {
		return false
	}

	if _, err := ds.s3.Get(repoFile.Name); err != nil {
		return false
	}

	return true
}

type LocalStorage struct {
	log          *zap.Logger
	videoPath    string
	chunkPath    string
	manifestPath string
}

func NewLocalStorage(log *zap.Logger, videoPath, chunkPath, manifestPath string) *LocalStorage {
	return &LocalStorage{
		log:          log,
		videoPath:    videoPath,
		chunkPath:    chunkPath,
		manifestPath: manifestPath,
	}
}

func (ls *LocalStorage) Store(video, manifest schema.InFile, chunks []schema.InFile) error {
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
		ls.log.Error(err.Error())
		return false
	}

	if _, err := os.Stat(filePath); err == nil {
		return true
	}

	return false
}

func (ls *LocalStorage) determinePath(filename string) (filePath string, err error) {
	if strings.HasSuffix(filename, ".mp4") {
		filePath = fmt.Sprintf("%v/%v", ls.videoPath, filename)
	} else if strings.HasSuffix(filename, ".m3u8") {
		filePath = fmt.Sprintf("%v/%v", ls.manifestPath, filename)
	} else if strings.HasSuffix(filename, ".ts") {
		filePath = fmt.Sprintf("%v/%v", ls.chunkPath, filename)
	} else {
		return filePath, fmt.Errorf("unsupported file format")
	}

	return filePath, nil
}
