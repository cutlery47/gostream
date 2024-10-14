package app

import (
	"log"

	"github.com/cutlery47/gostream/config"
	"github.com/cutlery47/gostream/internal/controller"
	"github.com/cutlery47/gostream/internal/service"
	"github.com/cutlery47/gostream/internal/storage"
	"github.com/cutlery47/gostream/internal/storage/obj"
	"github.com/cutlery47/gostream/internal/storage/repo"
	"github.com/cutlery47/gostream/pkg/logger"
	"github.com/cutlery47/gostream/pkg/server"
)

func Run() {
	config, err := config.New()
	if err != nil {
		log.Fatal("error when loading config:", err)
	}

	requestLogger := logger.New(config.Log.RequestLogsPath, false)
	errLogger := logger.New(config.Log.ErrorLogsPath, true)
	infoLogger := logger.New(config.Log.InfoLogsPath, false)

	// flushing any remaining data
	defer requestLogger.Sync()
	defer errLogger.Sync()
	defer infoLogger.Sync()

	if errLogger == nil || requestLogger == nil || infoLogger == nil {
		log.Fatal("logger paths should be fully provided")
	}

	var manifestStorage storage.Storage
	var chunkStorage storage.Storage
	var videoStorage storage.Storage

	if config.Storage.StorageType == "local" {
		manifestStorage = storage.NewLocalManifestStorage(config.Storage.Local.ManifestPath)
		chunkStorage = storage.NewLocalChunkStorage(config.Storage.Local.ChunkPath)
		videoStorage = storage.NewLocalVideoStorage(config.Storage.Local.VideoPath)
	} else {
		fileRepository, err := repo.NewFileRepository(config.Storage.Distr.DBConfig)
		if err != nil {
			log.Fatal(err)
		}

		s3, err := obj.NewS3(config.Storage.Distr.S3Config)
		if err != nil {
			log.Fatal(err)
		}

		manifestStorage = storage.NewDistributedManifestStorage(fileRepository, s3)
		return
	}

	chunkHandler := service.NewChunkHandler(infoLogger, chunkStorage)
	manifestHandler := service.NewManifestHandler(infoLogger, manifestStorage)
	videoHandler := service.NewVideoHandler(infoLogger, videoStorage)

	service := service.NewStreamService(
		chunkHandler,
		manifestHandler,
		videoHandler,
		infoLogger,
	)

	controller := controller.New(
		service,
		requestLogger,
		errLogger,
		infoLogger,
	)

	server := server.New(controller.Handler())

	server.Run()
}
