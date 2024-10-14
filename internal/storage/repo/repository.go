package repo

import (
	"database/sql"
	"fmt"

	"github.com/cutlery47/gostream/config"
	_ "github.com/lib/pq"
)

type CreateReadDeleteRepository interface {
	Create(file InRepositoryFile) error
	Get(filename string) (*RepositoryFile, error)
	Delete(filename string) error
}

type FileRepository struct {
	CreateReadDeleteRepository
	db *sql.DB
}

func (fr FileRepository) Create(file InRepositoryFile) error {
	return fmt.Errorf("xyu")
}

func (fr FileRepository) Get(filename string) (*RepositoryFile, error) {
	return nil, fmt.Errorf("xyu3")
}

func (fr FileRepository) Delete(filename string) error {
	return fmt.Errorf("xyu2")
}

func NewFileRepository(conf config.DBConfig) (*FileRepository, error) {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error when connecting to db: %v", err)
	}

	return &FileRepository{db: db}, nil
}

type ManifestRepository struct {
	fileRepo *FileRepository
}

func NewManifestRepository(conf config.DBConfig) (*ManifestRepository, error) {
	fileRepo, err := NewFileRepository(conf)
	if err != nil {
		return nil, err
	}

	return &ManifestRepository{
		fileRepo: fileRepo,
	}, nil
}

type VideoRepository struct {
	fileRepo *FileRepository
}

func NewVideoRepository(conf config.DBConfig) (*VideoRepository, error) {
	fileRepo, err := NewFileRepository(conf)
	if err != nil {
		return nil, err
	}

	return &VideoRepository{
		fileRepo: fileRepo,
	}, nil
}

type ChunkRepository struct {
	fileRepo *FileRepository
}

func NewChunkRepository(conf config.DBConfig) (*ChunkRepository, error) {
	fileRepo, err := NewFileRepository(conf)
	if err != nil {
		return nil, err
	}

	return &ChunkRepository{
		fileRepo: fileRepo,
	}, nil
}
