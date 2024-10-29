package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/cutlery47/gostream/config"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
		// check if video name is unique
		if pgerr, ok := err.(*pq.Error); ok && pgerr.Code == "23505" {
			return ErrUniqueVideo
		}
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

	if err := fr.insertFile(tx, id, video.VideoName, video.File.Size); err != nil {
		return err
	}

	if err := fr.insertLocation(tx, id, video.VideoName, video.File.Location); err != nil {
		return err
	}

	return nil
}

func (fr *FileRepository) createManifest(tx *sql.Tx, manifest InFile) error {
	id := uuid.New()

	if err := fr.insertFile(tx, id, manifest.Name, manifest.Size); err != nil {
		return err
	}

	if err := fr.insertLocation(tx, id, manifest.Name, manifest.Location); err != nil {
		return err
	}

	return nil

}

func (fr *FileRepository) createChunks(tx *sql.Tx, chunks []InFile) error {
	for _, chunk := range chunks {
		id := uuid.New()

		if err := fr.insertFile(tx, id, chunk.Name, chunk.Size); err != nil {
			return err
		}

		if err := fr.insertLocation(tx, id, chunk.Name, chunk.Location); err != nil {
			return err
		}
	}

	return nil
}

func (fr *FileRepository) insertFile(tx *sql.Tx, id uuid.UUID, filename string, size int) error {
	insertFile :=
		`
		INSERT INTO file_schema.files
		(id, filename)
		VALUES
		($1, $2);
		`

	insertMeta :=
		`
		INSERT INTO file_schema.files_meta
		(size, file_id)
		VALUES
		($1, $2);
		`

	if _, err := tx.Exec(insertFile, id, filename); err != nil {
		return err
	}

	if _, err := tx.Exec(insertMeta, size, id); err != nil {
		return err
	}

	return nil
}

func (fr *FileRepository) insertLocation(tx *sql.Tx, id uuid.UUID, filename string, location string) error {
	insert := fmt.Sprintf(
		`
		INSERT INTO file_schema.%v
		(location, file_id)
		VALUES
		($1, $2);
		`,
		fr.determineTable(filename),
	)

	if _, err := tx.Exec(insert, location, id); err != nil {
		return err
	}

	return nil
}

func (fr *FileRepository) determineTable(filename string) string {
	if strings.HasSuffix(filename, ".m3u8") {
		return "files_manifest"
	}

	if strings.HasSuffix(filename, ".ts") {
		return "files_ts"
	}

	return "files_vid"
}
