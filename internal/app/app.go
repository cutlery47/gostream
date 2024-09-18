package app

import (
	"github.com/cutlery47/gostream/internal/controller"
	"github.com/cutlery47/gostream/internal/service"
	"github.com/cutlery47/gostream/pkg/server"
)

func Run() {
	controller := controller.New(
		service.NewManifestService(),
		service.NewChunkService(),
	)

	server := server.New(controller.Handler())

	server.Run()
}
