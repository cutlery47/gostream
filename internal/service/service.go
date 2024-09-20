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
	Upload(file *os.File) error
}

// TODO: figure out how to isolate manifestService from chunkStorage

type manifestService struct {
	log          *zap.Logger
	errHander    errHandler
	manStorage   storage.Storage
	vidStorage   storage.Storage
	chunkStorage storage.Storage // trash
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
		ms.log.Info(fmt.Sprintf("Manifest for file (%v) is stored!", filename))
		return manifest, nil
	}

	// check if video file is even stored
	if !ms.vidStorage.Exists(fmt.Sprintf("%v.mp4", filename)) {
		return nil, ErrVideoNotFound
	}

	return ms.createManifest(filename)
}

func (ms *manifestService) createManifest(filename string) (*os.File, error) {
	chunkDir := ms.chunkStorage.Path()
	manDir := ms.manStorage.Path()
	vidDir := ms.vidStorage.Path()

	ms.checkOrCreateDirs(chunkDir, manDir, filename)

	cmd := utils.SegmentVideoAndCreateManifest(
		// precise file path
		fmt.Sprintf("%v/%v.mp4", vidDir, filename),
		// precise manifest path
		fmt.Sprintf("%v/%v.m3u8", manDir, filename),
		// chunk file path + template for segmentation
		fmt.Sprintf("%v/%v/%v_%%4d.ts", chunkDir, filename, filename),
		ms.segTime,
	)

	// check if segmentation went smoothely
	_, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// newly created manifest retrieval
	manifest, err := ms.manStorage.Get(fmt.Sprintf("%v.m3u8", filename))
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func (ms *manifestService) checkOrCreateDirs(chunkDir, manDir, filename string) {
	// creating necessary directories, if nonexistant
	utils.MKDir(chunkDir).Run()                                 // chunk dir
	utils.MKDir(fmt.Sprintf("%v/%v", chunkDir, filename)).Run() // chunk file dir
	utils.MKDir(manDir).Run()                                   // manifest dir
}

func (ms *manifestService) Upload(file *os.File) error {
	return ErrNotImplemented
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
	chunk, err := cs.storage.Get(filename)
	if err != nil {
		return nil, ErrChunkNotFound
	}

	return chunk, nil
}

func (cs *chunkService) Upload(file *os.File) error {
	return ErrNotImplemented
}

type videoService struct {
}
