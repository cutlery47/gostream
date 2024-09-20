package app

import (
	"log"

	"github.com/cutlery47/gostream/config"
	"github.com/cutlery47/gostream/internal/controller"
	"github.com/cutlery47/gostream/internal/service"
	"github.com/cutlery47/gostream/internal/storage"
	"github.com/cutlery47/gostream/pkg/logger"
	"github.com/cutlery47/gostream/pkg/server"
)

func Run() {
	config, err := config.New()
	if err != nil {
		log.Printf("error when loading config: %v", err)
		return
	}

	requestLogger := logger.New(config.Log.RequestLogsPath, false)
	errLogger := logger.New(config.Log.ErrorLogsPath, true)
	infoLogger := logger.New(config.Log.InfoLogsPath, false)

	// flushing any remaining data
	defer requestLogger.Sync()
	defer errLogger.Sync()
	defer infoLogger.Sync()

	if errLogger == nil || requestLogger == nil || infoLogger == nil {
		log.Println("logger paths should be fully provided")
		return
	}

	var manifestStorage storage.Storage
	var chunkStorage storage.Storage
	var videoStorage storage.Storage

	if config.Storage.StorageType == "local" {
		manifestStorage = storage.NewLocalManifestStorage(config.Storage.Local.ManifestPath)
		chunkStorage = storage.NewLocalChunkStorage(config.Storage.Local.ChunkPath)
		videoStorage = storage.NewLocalVideoStorage(config.Storage.Local.VideoPath)
	} else {
		log.Println("the only implemented storage type is local (so far)")
		return
	}

	manifestService := service.NewManifestService(
		infoLogger,
		errLogger,
		manifestStorage,
		videoStorage,
		chunkStorage,
		config.Segment.Time,
	)

	chunkService := service.NewChunkService(
		infoLogger,
		errLogger,
		chunkStorage,
	)

	videoService := service.NewVideoService(
		infoLogger,
		errLogger,
		videoStorage,
	)

	controller := controller.New(
		chunkService,
		manifestService,
		videoService,
		requestLogger,
		errLogger,
		infoLogger,
	)

	server := server.New(controller.Handler())

	server.Run()
}
