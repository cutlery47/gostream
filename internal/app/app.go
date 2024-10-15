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

	var store storage.Storage

	if config.Storage.StorageType == "local" {
		store = storage.NewLocalStorage(
			errLogger,
			config.Storage.Local.VideoPath,
			config.Storage.Local.ChunkPath,
			config.Storage.Local.ManifestPath,
		)
	} else {
		repo, err := repo.NewFileRepository(config.Storage.Distr.DBConfig)
		if err != nil {
			log.Fatal(err)
		}

		s3, err := obj.NewS3(config.Storage.Distr.S3Config)
		if err != nil {
			log.Fatal(err)
		}

		store = storage.NewDistibutedStorage(infoLogger, repo, s3)
	}

	manifestService := service.NewManifestService(infoLogger)

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
