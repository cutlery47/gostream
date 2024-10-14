package repo

import (
	"database/sql"
)

type Repository interface {
	Create(file InRepositoryFile) error
	Get(filename string) RepositoryFile
	Delete(filename string) bool
}

type ManifestRepository struct {
	db sql.DB
}

func NewManifestRepository() *ManifestRepository {
	return &ManifestRepository{}
}
