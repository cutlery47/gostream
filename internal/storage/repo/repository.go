package repo

import (
	"database/sql"
	"fmt"

	"github.com/cutlery47/gostream/config"
	_ "github.com/lib/pq"
)

type Repository interface {
	CreateAll(video InVideo, manifest InFile, chunks []InFile) error
	Read(filename string) (*File, error)
	Delete(filename string) (*File, error)
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

func (fr *FileRepository) CreateAll(video InVideo, manifest InFile, chunks []InFile) error {
	return fmt.Errorf("123123")
}

func (fr *FileRepository) Read(filename string) (*File, error) {
	return nil, fmt.Errorf("xyu3")
}

func (fr *FileRepository) Delete(filename string) (*File, error) {
	return nil, fmt.Errorf("xyu2")
}
