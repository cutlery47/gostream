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

	requestLogger := logger.New(config.Log.AppLogsPath+"/request.log", false)
	errLogger := logger.New(config.Log.AppLogsPath+"/error.log", true)
	infoLogger := logger.New(config.Log.AppLogsPath+"/info.log", false)

	// flushing any remaining data
	defer requestLogger.Sync()
	defer errLogger.Sync()
	defer infoLogger.Sync()

	if errLogger == nil || requestLogger == nil || infoLogger == nil {
		log.Fatal("all loggers should be properly configured")
	}

	var store storage.Storage

	paths := storage.Paths{
		VidPath:   config.Storage.Local.VideoPath,
		ManPath:   config.Storage.Local.ManifestPath,
		ChunkPath: config.Storage.Local.ChunkPath,
	}

	if config.Storage.StorageType == "local" {
		store = storage.NewLocalStorage(errLogger, paths)
	} else {
		repo, err := repo.NewFileRepository(config.Storage.Distr.DBConfig)
		if err != nil {
			log.Fatal("Error when initializing db: ", err)
		}

		s3, err := obj.NewS3(config.Storage.Distr.S3Config)
		if err != nil {
			log.Fatal("Error when initializing s3: ", err)
		}

		store = storage.NewDistibutedStorage(infoLogger, errLogger, paths, repo, s3)
	}

	manifestService := service.NewManifestService(infoLogger, errLogger)

	service := service.NewStreamService(
		infoLogger,
		store,
		manifestService,
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
