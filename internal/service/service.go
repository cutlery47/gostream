package service

import (
	"fmt"
	"io"

	"github.com/cutlery47/gostream/internal/utils"
	"go.uber.org/zap"
)

type Service interface {
	ServeManifest(filename string) (io.Reader, error)
	ServeChunk(filename string) (io.Reader, error)
	ServeEntireVideo(filename string) (io.Reader, error)
	UploadVideo(file io.Reader, filename string) error
	RemoveVideo(filename string) error
}

type StreamService struct {
	chunkHandler    RemoveRetriever
	manifestHandler RemoveRetriever
	videoHandler    UploadRemoveRetriever
	log             *zap.Logger
}

func NewStreamService(
	chunkHandler RemoveRetriever,
	manifestHandler RemoveRetriever,
	videoHandler UploadRemoveRetriever,
	log *zap.Logger,
) *StreamService {
	return &StreamService{
		chunkHandler:    chunkHandler,
		manifestHandler: manifestHandler,
		videoHandler:    videoHandler,
		log:             log,
	}
}

// TODO: add chunk checks
func (ss *StreamService) ServeManifest(filename string) (io.Reader, error) {
	// check if we already store the manifest file
	if manifest, _ := ss.manifestHandler.Retrieve(filename); manifest != nil {
		ss.log.Info("Manifest for %v is already stored. Returning...")
		return manifest, nil
	}

	// check if video file is even stored
	if !ss.videoHandler.Exists(fmt.Sprintf("%v.mp4", filename)) {
		return nil, ErrVideoNotFound
	}

	chunkDir := ss.chunkHandler.Path()
	chunkFileDir := fmt.Sprintf("%v/%v", ss.chunkHandler.Path(), filename)
	manifestDir := ss.manifestHandler.Path()
	videoDir := ss.videoHandler.Path()

	ss.createDirs(chunkDir, chunkFileDir, manifestDir, videoDir)

	cmd := utils.SegmentVideoAndCreateManifest(
		// precise file path
		fmt.Sprintf("%v/%v.mp4", videoDir, filename),
		// precise manifest path
		fmt.Sprintf("%v/%v.m3u8", manifestDir, filename),
		// chunk file path + template for segmentation
		fmt.Sprintf("%v/%v_%%4d.ts", chunkFileDir, filename),
	)

	// check if segmentation went smoothely
	if _, err := cmd.Output(); err != nil {
		return nil, err
	}

	return ss.manifestHandler.Retrieve(filename)
}

func (ss *StreamService) createDirs(chunkDir, chunkFileDir, manifestDir, videoDir string) {
	utils.MKDir(chunkDir)
	utils.MKDir(chunkFileDir)
	utils.MKDir(manifestDir)
	utils.MKDir(videoDir)
}

func (ss *StreamService) ServeChunk(filename string) (io.Reader, error) {
	if !ss.chunkHandler.Exists(filename) {
		return nil, ErrChunkNotFound
	}
	return ss.chunkHandler.Retrieve(filename)
}

func (ss *StreamService) ServeEntireVideo(filename string) (io.Reader, error) {
	if !ss.videoHandler.Exists(filename) {
		return nil, ErrVideoNotFound
	}
	return ss.videoHandler.Retrieve(filename)
}

func (ss *StreamService) UploadVideo(file io.Reader, filename string) error {
	return ss.videoHandler.Upload(file, filename)
}

// TODO: figure out how to make this method atomic
func (ss *StreamService) RemoveVideo(filename string) error {
	if err := ss.videoHandler.Remove(filename); err != nil {
		return err
	}
	if err := ss.manifestHandler.Remove(filename); err != nil {
		return err
	}
	if err := ss.chunkHandler.Remove(filename); err != nil {
		return err
	}
	return nil
}
