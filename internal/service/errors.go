package service

import (
	"errors"

	"go.uber.org/zap"
)

var (
	ErrManifestNotFound = newServiceError("couldn't find requested manifest file")
)

type ServiceError struct {
	err error
}

func newServiceError(message string) *ServiceError {
	return &ServiceError{
		err: errors.New(message),
	}
}

func (se ServiceError) Error() string { return se.err.Error() }

type errHandler struct {
	log *zap.Logger
}

func (eh errHandler) Handle(err error) {

}
