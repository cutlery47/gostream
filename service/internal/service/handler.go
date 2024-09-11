package service

import (
	"errors"
	"fmt"
	"gostream/service/pkg/server"
	"log"
	"net/http"
	"os"
)

type Handler interface {
	getVideoName() (name string, err error)
}

type HTTPHandler struct {
	ch chan string
}

func NewHTTPHandler() Handler {
	handler := &HTTPHandler{}

	serv := server.New(
		handler,
	)

	serv.Run()

	return handler
}

func (hh *HTTPHandler) getVideoName() (name string, err error) {
	name = <-hh.ch
	if name == "" {
		return "", fmt.Errorf("video name was not provided")
	}

	return name, nil
}

func (hh *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Body)
}

type LocalHandler struct {
}

func (lh *LocalHandler) getVideoName() (name string, err error) {
	if len(os.Args) < 2 {
		return "", errors.New("video name is required")
	}

	return os.Args[1], nil
}

type ErrHandler interface {
	Handle(err error)
}

type LocalErrHandler struct{}

func (leh *LocalErrHandler) Handle(err error) {
	log.Println(err)
}
