package service

import (
	"fmt"
	"os"

	"github.com/cutlery47/gostream/internal/storage"
	"github.com/cutlery47/gostream/internal/utils"
	"go.uber.org/zap"
)

type Service interface {
	Serve(filename string) (*os.File, error)
}

// TODO: figure out how to isolate manifestService from chunkStorage

type manifestService struct {
	log          *zap.Logger
	errHander    errHandler
	manStorage   storage.Storage
	vidStorage   storage.Storage
	chunkStorage storage.Storage
	// length of each segmented piece
	segTime int
}

func NewManifestService(infoLog, errLog *zap.Logger, manStorage, vidStorage, chunkStorage storage.Storage, segTime int) *manifestService {
	errHandler := errHandler{
		log: errLog,
	}

	return &manifestService{
		log:          infoLog,
		errHander:    errHandler,
		manStorage:   manStorage,
		vidStorage:   vidStorage,
		chunkStorage: chunkStorage,
		segTime:      segTime,
	}
}

func (ms *manifestService) Serve(filename string) (*os.File, error) {
	// check if we already store the manifest file
	manifest, err := ms.manStorage.Get(fmt.Sprintf("%v.m3u8", filename))
	if err != nil {
		ms.log.Info(fmt.Sprintf("Manifest for file (%v) is not stored... Creating", filename))
	} else {
		return manifest, nil
	}

	// check if video file is even stored
	if !ms.vidStorage.Exists(fmt.Sprintf("%v.mp4", filename)) {
		return nil, ErrVideoNotFound
	}

	out := utils.SegmentVideoAndCreateManifest(
		ms.vidStorage.Path(),
		ms.manStorage.Path(),
		ms.chunkStorage.Path(),
		ms.segTime,
	)
	ms.log.Info(out.String())

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
