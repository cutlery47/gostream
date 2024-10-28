package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cutlery47/gostream/config"
	"github.com/google/uuid"
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
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", conf.User, conf.Password, conf.Host, conf.Port, conf.DBName, conf.SSLMode)
	// openning db connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error when connecting to db: %v", err)
	}

	return &FileRepository{db: db}, nil
}

func (fr *FileRepository) CreateAll(video InVideo, manifest InFile, chunks []InFile) error {
	ctx := context.Background()

	tx, err := fr.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fr.createVideo(tx, video); err != nil {
		return err
	}

	if err := fr.createManifest(tx, manifest); err != nil {
		return err
	}

	if err := fr.createChunks(tx, chunks); err != nil {
		return err
	}

	return tx.Commit()
}

func (fr *FileRepository) Read(filename string) (*File, error) {
	return nil, fmt.Errorf("xyu3")
}

func (fr *FileRepository) Delete(filename string) (*File, error) {
	return nil, fmt.Errorf("xyu2")
}

func (fr *FileRepository) createVideo(tx *sql.Tx, video InVideo) error {
	id := uuid.New()

	preparedInsertVideoFile :=
		`
		INSERT INTO file_schema.files
		(id, filename)
		VALUES
		($1, $2);
		`

	insertVideoFile, err := tx.Prepare(preparedInsertVideoFile)
	if err != nil {
		return err
	}

	preparedInsertVideoMeta :=
		`
		INSERT INTO file_schema.files_meta
		(size, file_id)
		VALUES
		($1, $2);
		`

	insertVideoMeta, err := tx.Prepare(preparedInsertVideoMeta)
	if err != nil {
		return err
	}

	preparedInsertVideo :=
		`
		INSERT INTO file_schema.files_vid
		(location, file_id)
		VALUES
		($1, $2);
		`

	insertVideo, err := tx.Prepare(preparedInsertVideo)
	if err != nil {
		return err
	}

	if _, err := insertVideoFile.Exec(id, video.VideoName); err != nil {
		return err
	}

	if _, err := insertVideoMeta.Exec(video.File.Size, id); err != nil {
		return err
	}

	if _, err := insertVideo.Exec(video.File.Location, id); err != nil {
		return err
	}

	return nil
}

func (fr *FileRepository) createManifest(tx *sql.Tx, manifest InFile) error {
	return nil
}

func (fr *FileRepository) createChunks(tx *sql.Tx, chunks []InFile) error {
	return nil
}
