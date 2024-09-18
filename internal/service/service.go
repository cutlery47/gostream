package service

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"
// )

// const (
// 	// directory for storing videos
// 	videoDir = "vids"
// 	// directory for storing video segments
// 	segmentDir = videoDir + "/segmented"
// 	// length (in seconds) of a single video segment
// 	segmentTime = 4
// )

// type httpHandler struct {
// 	// request channel
// 	reqCh chan string
// 	// response channel
// 	resCh chan httpResponse
// 	// exit channel
// 	exitCh chan bool
// }

// type httpResponse struct {
// 	// error, occurred when handling request
// 	err error
// 	// response file
// 	file []byte
// }

// func (hh *httpHandler) getMediaName() (name string, err error) {
// 	// receiving request from http server
// 	name = <-hh.reqCh
// 	if name == "" {
// 		// if name is empty -- we assume it hasn't been provided
// 		err = fmt.Errorf("video name was not provided")
// 		hh.sendResponse(nil, err)
// 		return "", err
// 	}

// 	return name, nil
// }

// func (hh *httpHandler) serveMedia(fileName string) (err error) {
// 	// check which type of file was requested
// 	if strings.HasSuffix(fileName, ".ts") {
// 		// stream chunk was requested
// 		return hh.serveChunk(fileName)
// 	} else {
// 		// manifest file was requested
// 		return hh.serveManifest(fileName)
// 	}
// }

// func (hh *httpHandler) serveChunk(chunkName string) (err error) {
// 	chunkDir := removeSuffix(chunkName, "_")
// 	// retrieving requested chunk
// 	chunk := hh.getChunk(chunkName, chunkDir)
// 	if err != nil {
// 		hh.sendResponse(nil, err)
// 		return err
// 	}
// 	// responding with the requested chunk
// 	hh.sendFileResponse(chunk)

// 	return nil
// }

// func (hh *httpHandler) serveManifest(videoName string) (err error) {
// 	// checking if manifest file already exists, return one if it does
// 	manifest := hh.getManifest(videoName)
// 	if manifest != nil {
// 		hh.sendFileResponse(manifest)
// 		return nil
// 	}

// 	// checking if requested video exists
// 	cmd := findmedia(videoName)
// 	log.Println(cmd.Output())

// 	cmd = segment(videoName)
// 	log.Println("segment: ", cmd.String())

// 	// responding
// 	manifest = hh.getManifest(videoName)
// 	hh.sendFileResponse(manifest)

// 	return nil
// }

// func (hh *httpHandler) getManifest(videoName string) *os.File {
// 	manifest, err := os.Open(fmt.Sprintf("%v/%v/%v.m3u8", segmentDir, videoName, videoName))
// 	if err != nil {
// 		log.Println("getManifest:", err)
// 		return nil
// 	}

// 	return manifest
// }

// func (hh *httpHandler) getChunk(chunkName, chunkDir string) *os.File {
// 	chunk, err := os.Open(fmt.Sprintf("%v/%v/%v", segmentDir, chunkDir, chunkName))
// 	if err != nil {
// 		log.Println("getChunk:", err)
// 		return nil
// 	}

// 	return chunk
// }

// func (hh *httpHandler) sendFileResponse(file *os.File) error {
// 	// retrieve file metadata
// 	meta, err := file.Stat()
// 	if err != nil {
// 		hh.sendResponse(nil, err)
// 		return err
// 	}

// 	// buffer for storing file in memory
// 	fileBuffer := make([]byte, meta.Size())
// 	// reading file data into the buffer
// 	if _, err = bufio.NewReader(file).Read(fileBuffer); err != nil {
// 		hh.sendResponse(nil, err)
// 		return err
// 	}

// 	// sending the file
// 	hh.sendResponse(fileBuffer, nil)

// 	return nil
// }

// func (hh *httpHandler) sendResponse(file []byte, err error) {
// 	// sending error message as a response
// 	hh.resCh <- httpResponse{
// 		err:  err,
// 		file: file,
// 	}
// }

// func (hh *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")

// 	//
// 	go hh.serve(w)

// 	fileName := r.RequestURI[1:]
// 	// sending the requested media name over to the handler
// 	hh.reqCh <- fileName
// 	// waiting for handler to finish
// 	<-hh.exitCh
// }

// func (hh *httpHandler) serve(w http.ResponseWriter) {
// 	// returning response from the handler
// 	response := <-hh.resCh

// 	if response.err != nil {
// 		_, err := w.Write([]byte(response.err.Error()))
// 		fmt.Println(err)
// 	} else {
// 		_, err := w.Write(response.file)
// 		fmt.Println(err)
// 	}

// 	// signaling that the response has been sent
// 	hh.exitCh <- true
//
