package app

import (
	"log"

	"github.com/cutlery47/gostream/config"
	v1 "github.com/cutlery47/gostream/internal/controller/http/v1"
	"github.com/cutlery47/gostream/internal/service"
	"github.com/cutlery47/gostream/internal/storage"
	"github.com/cutlery47/gostream/pkg/httpserver"
	"github.com/cutlery47/gostream/pkg/logger"
	"github.com/labstack/echo/v4"
)

//	@title			Gostream
//	@version		1.0
//	@description	A simple golang streaming service.

//	@contact.name	Arkhip Ivanchenko
//	@contact.url	https://github.com/cutlery47
//	@contact.email	kitchen_cutlery@mail.ru

//	@host	localhost:8080

func Run() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("error when loading config: ", err)
	}

	log.Println(cfg)

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

	if cfg.Flag.Type == "local" {
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

		st = storage.NewDistibutedStorage(infLog, errLog, cfg.Storage, repo, s3)
	}

	svc := service.NewStreamService(
		infLog,
		cfg.Storage.Local,
		st,
	)

	e := echo.New()
	v1.NewController(e, svc, reqLog, errLog, infLog)

	httpserver.New(e).Run()
}
