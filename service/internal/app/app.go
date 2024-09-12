package app

import (
	"gostream/service/internal/service"
)

func Run() {
	handler := service.NewHttpHandler()

	service := service.New(handler)
	service.Run()
}
