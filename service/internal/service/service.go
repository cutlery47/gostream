package service

import (
	"log"
)

type Service struct {
	handler Handler
}

func (s *Service) handleErr(err error) {
	log.Println(err)
}

func (s *Service) Run() {
	for {
		fileName, err := s.handler.getMediaName()
		if err != nil {
			s.handleErr(err)
		}

		err = s.handler.serveMedia(fileName)
		if err != nil {
			s.handleErr(err)
		}
	}
}

func New(handler Handler) *Service {
	return &Service{
		handler: handler,
	}
}
