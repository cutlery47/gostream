package app

import (
	"gostream/service/internal/service"
)

func Run() {
	eHandler := &service.LocalErrHandler{}
	handler := &service.LocalHandler{}

	service := service.New(handler, eHandler)
	service.Run()
}
