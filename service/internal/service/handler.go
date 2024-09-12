package service

import (
	"bufio"
	"errors"
	"fmt"
	"gostream/service/pkg/server"
	"log"
	"net/http"
	"os"
)

const videoDir = "vids"

type Handler interface {
	// retrieve the name of the video to be streamed
	getVideoName() (name string, err error)

	// find the video on service disc
	findVideo(videoName string) (err error)

	// find or create a directory for storing segmented files
	findOrCreateDir(videoName string)

	// split the initial video into slices
	segmentVideo(videoName string) (err error)

	// return the manifest file for download
	returnManifest(videoName string) (err error)
}

type httpHandler struct {
	// request channel
	reqCh chan string

	// response channel
	resCh chan httpResponse

	// exit channel
	exitCh chan bool
}

type httpResponse struct {
	err      error
	manifest []byte
}

func NewHttpHandler() *httpHandler {
	handler := &httpHandler{
		reqCh:  make(chan string),
		resCh:  make(chan httpResponse),
		exitCh: make(chan bool),
	}

	serv := server.New(
		handler,
	)

	go serv.Run()

	return handler
}

func (hh *httpHandler) getVideoName() (name string, err error) {
	// receiving request from http server
	name = <-hh.reqCh
	if name == "" {
		err = fmt.Errorf("video name was not provided")
		// sending error message as a response
		hh.resCh <- httpResponse{
			err:      err,
			manifest: nil,
		}
		// waiting until the error message has been sent
		<-hh.exitCh
		return "", err
	}

	return name, nil
}

func (hh *httpHandler) findVideo(videoName string) (err error) {
	// script for finding video file
	cmd := findmedia(videoName, videoDir)

	if _, err = cmd.Output(); err != nil {
		// sending error message as a response
		hh.resCh <- httpResponse{
			err:      err,
			manifest: nil,
		}
		// waiting until the error message has been sent
		<-hh.exitCh
		return err
	}

	return nil
}

func (hh *httpHandler) findOrCreateDir(videoName string) {
	// script for creating a directory if one doesnt exist
	cmd := mkdir(videoName, videoDir)
	cmd.Run()
}

func (hh *httpHandler) segmentVideo(videoName string) (err error) {
	// ffmpeg segmentation script
	cmd := segment(videoName, videoDir)

	if _, err := cmd.Output(); err != nil {
		// sending error message as a response
		hh.resCh <- httpResponse{
			err:      err,
			manifest: nil,
		}
		// waiting until the error message has been sent
		<-hh.exitCh
		return fmt.Errorf("segmentation error: %v", err)
	}

	return nil
}

func (hh *httpHandler) returnManifest(videoName string) (err error) {
	file, err := os.Open(fmt.Sprintf("%v/%v.m3u8", videoDir, videoName))
	defer file.Close()
	if err != nil {
		hh.resCh <- httpResponse{
			err:      err,
			manifest: nil,
		}
		<-hh.exitCh
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		hh.resCh <- httpResponse{
			err:      err,
			manifest: nil,
		}
		<-hh.exitCh
		return err
	}

	fileBytes := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(fileBytes)
	if err != nil {
		hh.resCh <- httpResponse{
			err:      err,
			manifest: nil,
		}
		<-hh.exitCh
		return err
	}

	hh.resCh <- httpResponse{
		err:      nil,
		manifest: fileBytes,
	}
	<-hh.exitCh

	return nil

}

func (hh *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// getting video name from the url string
	hh.reqCh <- r.RequestURI[1:]

	// returning response from the handler
	response := <-hh.resCh
	if response.err != nil {
		w.Write([]byte(response.err.Error()))
	} else {
		w.Header().Set("Content-Disposition", "attachment")
		w.Write(response.manifest)
	}

	hh.exitCh <- true
}

type localHandler struct{}

func NewLocalHandler() *localHandler {
	return &localHandler{}
}

func (lh *localHandler) getVideoName() (name string, err error) {
	if len(os.Args) < 2 {
		return "", errors.New("video name is required")
	}

	return os.Args[1], nil
}

func (lh *localHandler) findVideo(videoName string) (err error) {
	// script for finding video file
	cmd := findmedia(videoName, videoDir)

	if _, err = cmd.Output(); err != nil {
		return err
	}

	return nil
}

func (lh *localHandler) findOrCreateDir(videoName string) {
	// script for creating a directory if one doesnt exist
	cmd := mkdir(videoName, videoDir)
	cmd.Run()
}

func (lh *localHandler) segmentVideo(videoName string) (err error) {
	// ffmpeg segmentation script
	cmd := segment(videoName, videoDir)

	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("segmentation error: %v", err)
	}

	return nil
}

type ErrHandler interface {
	Handle(err error)
}

type localErrHandler struct{}

func (leh *localErrHandler) Handle(err error) {
	log.Println(err)
}

func NewLocalErrHandler() *localErrHandler {
	return &localErrHandler{}
}
