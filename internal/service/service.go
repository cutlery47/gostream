package service

import (
	"os"

	"github.com/cutlery47/gostream/internal/storage"
	"go.uber.org/zap"
)

const (
	// directory for storing videos
	videoDir = "vids"
	// directory for storing video segments
	segmentDir = videoDir + "/segmented"
	// length (in seconds) of a single video segment
	segmentTime = 4
)

type Service interface {
	Serve(filename string) (*os.File, error)
}

type manifestService struct {
	log       *zap.Logger
	errHander errHandler
	storage   storage.Storage
}

func NewManifestService(infoLog, errLog *zap.Logger, storage storage.Storage) *manifestService {
	errHandler := errHandler{
		log: errLog,
	}

	return &manifestService{
		log:       infoLog,
		errHander: errHandler,
		storage:   storage,
	}
}

func (ms *manifestService) Serve(filename string) (*os.File, error) {
	manifest, err := ms.storage.Get(filename)
	if err != nil {
		ms.log.Info("Requested manifest is not currently stored... Creating")
	}

	return manifest, nil
}

type chunkService struct {
	log        *zap.Logger
	errHandler errHandler
	storage    storage.Storage
}

func NewChunkService(infoLog, errLog *zap.Logger, storage storage.Storage) *chunkService {
	errHandler := errHandler{
		log: errLog,
	}

	return &chunkService{
		log:        infoLog,
		errHandler: errHandler,
		storage:    storage,
	}
}

func (cs *chunkService) Serve(filename string) (*os.File, error) {
	return nil, nil
}
