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
	cfg, err := config.New()
	if err != nil {
		log.Fatal("error when loading config:", err)
	}

	reqLog := logger.New(cfg.Log.AppLogsPath+"/request.log", false)
	errLog := logger.New(cfg.Log.AppLogsPath+"/error.log", true)
	infLog := logger.New(cfg.Log.AppLogsPath+"/info.log", false)

	// flushing any remaining data
	defer reqLog.Sync()
	defer errLog.Sync()
	defer infLog.Sync()

	if errLog == nil || reqLog == nil || infLog == nil {
		log.Fatal("all loggers should be properly configured")
	}

	var st storage.Storage

	if cfg.Storage.StorageType == "local" {
		st = storage.NewLocalStorage(errLog, cfg.Storage.Local)
	} else {
		repo, err := storage.NewFileRepository(cfg.Storage.Distr.DBConfig)
		if err != nil {
			log.Fatal("Error when initializing db: ", err)
		}

		s3, err := storage.NewS3(cfg.Storage.Distr.S3Config)
		if err != nil {
			log.Fatal("Error when initializing s3: ", err)
		}

		st = storage.NewDistibutedStorage(infLog, errLog, cfg.Storage.Distr, repo, s3)
	}

	svc := service.NewStreamService(
		infLog,
		cfg.Storage.Local,
		st,
	)

	ctr := controller.New(
		svc,
		reqLog,
		errLog,
		infLog,
	)

	srv := server.New(ctr.Handler())

	srv.Run()
}
