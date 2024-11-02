package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cutlery47/gostream/config"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository interface {
	// creates all the entries in db
	CreateAll(ctx context.Context, video, manifest File, chunks []File) error
	// returns object storage location of a certain file
	Read(ctx context.Context, filename string) (Location, error)
	// deletes file from db and returns its object storage location
	Delete(ctx context.Context, filename string) (Location, error)
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

func (fr *FileRepository) CreateAll(ctx context.Context, video File, manifest File, chunks []File) error {
	tx, err := fr.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fr.insertFile(ctx, tx, video); err != nil {
		// check if video name is unique
		if pgerr, ok := err.(*pq.Error); ok && pgerr.Code == "23505" {
			err = ErrUniueVideo
		}
		return err
	}

	if err := fr.insertFile(ctx, tx, manifest); err != nil {
		return err
	}

	// TODO: try out goroutine based version

	// dumb af but aight
	for _, chunk := range chunks {
		if err := fr.insertFile(ctx, tx, chunk); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (fr *FileRepository) Read(ctx context.Context, filename string) (location Location, err error) {
	query :=
		`
		SELECT bucket, object 
		FROM file_schema.files
		WHERE name = $1
		`

	res := fr.db.QueryRowContext(ctx, query, filename)
	if err := res.Err(); err != nil {
		return location, err
	}

	if err := res.Scan(&location.Bucket, &location.Object); err != nil {
		return location, err
	}

	return location, err
}

func (fr *FileRepository) Delete(ctx context.Context, filename string) (location Location, err error) {
	query :=
		`
		DELETE file_schema.files AS f
		WHERE f.name = $1
		RETURNING f.bucket, f.object;
		`

	res := fr.db.QueryRowContext(ctx, query, filename)
	if err := res.Err(); err != nil {
		return location, err
	}

	if err := res.Scan(&location.Bucket, &location.Object); err != nil {
		return location, err
	}

	return location, err

}

func (fr *FileRepository) insertFile(ctx context.Context, tx *sql.Tx, file File) error {
	id := uuid.New()

	// query for file insertion
	insertFile :=
		`
		INSERT INTO file_schema.files
		(id, name, bucket, object)
		VALUES
		($1, $2, $3, $4);
		`

	// query for metadata insertion (along with file)
	insertMeta :=
		`
		INSERT INTO file_schema.files_meta
		(file_id)
		VALUES
		($1);
		`

	if _, err := tx.ExecContext(ctx, insertFile, id, file.FileName, file.Location.Bucket, file.Location.Object); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, insertMeta, id); err != nil {
		return err
	}

	return nil
}
