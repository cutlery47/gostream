package repo

import (
	"database/sql"
	"fmt"

	"github.com/cutlery47/gostream/config"
	_ "github.com/lib/pq"
)

type Repository interface {
	CreateAll(video, manifest InRepositoryFile, chunks []InRepositoryFile) error
	Read(filename string) (*RepositoryFile, error)
	Delete(filename string) (*RepositoryFile, error)
}

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(conf config.DBConfig) (*FileRepository, error) {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)
	// openning db connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error when connecting to db: %v", err)
	}

	return &FileRepository{db: db}, nil
}

func (fr *FileRepository) CreateAll(video, manifest InRepositoryFile, chunks []InRepositoryFile) error {
	return fmt.Errorf("123123")
}

func (fr *FileRepository) Read(filename string) (*RepositoryFile, error) {
	return nil, fmt.Errorf("xyu3")
}

func (fr *FileRepository) Delete(filename string) (*RepositoryFile, error) {
	return nil, fmt.Errorf("xyu2")
}

// type ManifestRepository struct {
// 	*FileRepository
// }

// func NewManifestRepository(conf config.DBConfig) (*ManifestRepository, error) {
// 	fileRepo, err := NewFileRepository(conf)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ManifestRepository{fileRepo}, nil
// }

// type VideoRepository struct {
// 	*FileRepository
// }

// func NewVideoRepository(conf config.DBConfig) (*VideoRepository, error) {
// 	fileRepo, err := NewFileRepository(conf)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &VideoRepository{fileRepo}, nil
// }

// type ChunkRepository struct {
// 	*FileRepository
// }

// func NewChunkRepository(conf config.DBConfig) (*ChunkRepository, error) {
// 	fileRepo, err := NewFileRepository(conf)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ChunkRepository{fileRepo}, nil
// }
