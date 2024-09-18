package app

import (
	"log"

	"github.com/cutlery47/gostream/config"
	"github.com/cutlery47/gostream/internal/controller"
	"github.com/cutlery47/gostream/internal/service"
	"github.com/cutlery47/gostream/pkg/logger"
	"github.com/cutlery47/gostream/pkg/server"
)

func Run() {
	config := config.New()

	requestLogger := logger.New(config.Log.RequestLogsPath, false)
	errLogger := logger.New(config.Log.ErrorLogsPath, true)
	infoLogger := logger.New(config.Log.InfoLogsPath, false)

	defer requestLogger.Sync()
	defer errLogger.Sync()
	defer infoLogger.Sync()

	if errLogger == nil || requestLogger == nil || infoLogger == nil {
		log.Println("logger paths should be fully provided")
		return
	}

	controller := controller.New(
		service.NewManifestService(),
		service.NewChunkService(),
		requestLogger,
		errLogger,
		infoLogger,
	)

	server := server.New(controller.Handler())

	server.Run()
}
