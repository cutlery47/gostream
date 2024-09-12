package app

import (
	"gostream/service/internal/service"
)

func Run() {
	eHandler := service.NewLocalErrHandler()
	handler := service.NewHttpHandler()

	service := service.New(handler, eHandler)
	service.Run()
}
