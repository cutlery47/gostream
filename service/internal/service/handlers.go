package service

import (
	"bufio"
	"errors"
	"fmt"
	"gostream/service/pkg/server"
	"net/http"
	"os"
	"strings"
)

const videoDir = "vids"

type Handler interface {
	// retrieve the name of the file to be streamed
	getMediaName() (fileName string, err error)

	// return the file to be streamed
	serveMedia(fileName string) (err error)
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
	// error, occurred when handling request
	err error

	// response file
	file []byte
}

func NewHttpHandler() *httpHandler {
	handler := &httpHandler{
		reqCh:  make(chan string),
		resCh:  make(chan httpResponse),
		exitCh: make(chan bool),
	}

	// creating a new server and listening to http requests
	serv := server.New(
		handler,
	)
	go serv.Run()

	return handler
}

func (hh *httpHandler) handleErr(err error) {
	// sending error message as a response
	hh.resCh <- httpResponse{
		err:  err,
		file: nil,
	}
	// waiting until the error message has been sent
	<-hh.exitCh
}

func (hh *httpHandler) getMediaName() (name string, err error) {
	// receiving request from http server
	name = <-hh.reqCh
	if name == "" {
		err = fmt.Errorf("video name was not provided")
		hh.handleErr(err)
		return "", err
	}

	return name, nil
}

func (hh *httpHandler) serveMedia(fileName string) (err error) {
	if strings.HasSuffix(fileName, ".ts") {
		// stream chunk was requested
		return hh.serveChunk(fileName)
	} else {
		// manifest file was requested
		return hh.serveManifest(fileName)
	}
}

func (hh *httpHandler) serveChunk(chunkName string) (err error) {
	chunkFile, err := os.Open(fmt.Sprintf("vids/segmented/media/%v", chunkName))
	if err != nil {
		hh.handleErr(err)
		return err
	}
	defer chunkFile.Close()

	stat, err := chunkFile.Stat()
	if err != nil {
		hh.handleErr(err)
		return err
	}

	fileBytes := make([]byte, stat.Size())
	bufio.NewReader(chunkFile).Read(fileBytes)

	// returning the file
	hh.resCh <- httpResponse{
		err:  nil,
		file: fileBytes,
	}

	return nil
}

func (hh *httpHandler) serveManifest(mediaName string) (err error) {
	// script for finding video file
	cmd := findmedia(mediaName, videoDir)
	if _, err = cmd.Output(); err != nil {
		hh.handleErr(fmt.Errorf("error when trying to find video"))
		return err
	}

	// script for creating a directory if one doesnt exist
	cmd = mkdir(mediaName, videoDir)
	cmd.Run()

	// ffmpeg segmentation script
	cmd = segment(mediaName, videoDir)
	if _, err := cmd.Output(); err != nil {
		hh.handleErr(fmt.Errorf("segmentation error"))
		return err
	}

	// getting manifest file descriptor
	file, err := os.Open(fmt.Sprintf("%v/segmented/%v/%v.m3u8", videoDir, mediaName, mediaName))
	if err != nil {
		hh.handleErr(err)
		return err
	}
	defer file.Close()

	// getting manifest metadata
	stat, err := file.Stat()
	if err != nil {
		hh.handleErr(err)
		return err
	}

	// converting the file to byte stream and sending to the user
	fileBytes := make([]byte, stat.Size())
	if _, err = bufio.NewReader(file).Read(fileBytes); err != nil {
		hh.handleErr(err)
		return err
	}

	// sending manifest back
	hh.resCh <- httpResponse{
		err:  nil,
		file: fileBytes,
	}
	// waiting until response is sent
	<-hh.exitCh

	return nil

}

func (hh *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	go hh.serve(w)

	hh.reqCh <- r.RequestURI[1:]
}

func (hh *httpHandler) serve(w http.ResponseWriter) {
	// returning response from the handler
	response := <-hh.resCh

	if response.err != nil {
		w.Write([]byte(response.err.Error()))
	} else {
		w.Write(response.file)
	}

	// signaling that the response has been sent
	hh.exitCh <- true
}

type localHandler struct{}

func NewLocalHandler() *localHandler {
	return &localHandler{}
}

func (lh *localHandler) getFileName() (name string, err error) {
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
