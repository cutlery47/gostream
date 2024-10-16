package storage

import "errors"

var (
	ErrNotImplemented        = newStorageError("feature is not yet implemented")
	ErrUnsupportedFileFormat = newStorageError("unsupported file format")
)

type StorageError struct {
	err error
}

func newStorageError(message string) *StorageError {
	return &StorageError{
		err: errors.New(message),
	}
}

func (se StorageError) Error() string { return se.err.Error() }
