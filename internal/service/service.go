package service

import (
	"fmt"

	"github.com/cutlery47/gostream/internal/schema"
	"github.com/cutlery47/gostream/internal/utils"
	"go.uber.org/zap"
)

type Service interface {
	ServeManifest(filename string) (*schema.OutFile, error)
	ServeChunk(filename string) (*schema.OutFile, error)
	ServeEntireVideo(filename string) (*schema.OutFile, error)
	UploadVideo(file schema.InFile) error
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
func (ss *StreamService) ServeManifest(filename string) (*schema.OutFile, error) {
	// check if we already store the manifest file
	if manifest, _ := ss.manifestHandler.Retrieve(filename + ".m3u8"); manifest != nil {
		ss.log.Info(fmt.Sprintf("Manifest for %v is already stored. Returning...", filename))
		return manifest, nil
	}

	// check if video file is even stored
	if !ss.videoHandler.Exists(fmt.Sprintf("%v.mp4", filename)) {
		return nil, ErrVideoNotFound
	}

	chunkDir := ss.chunkHandler.Path()
	chunkFileDir := fmt.Sprintf("%v/%v", chunkDir, filename)
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
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, ErrSegmentationException
	}

	ss.log.Error(string(out))

	return ss.manifestHandler.Retrieve(filename + ".m3u8")
}

func (ss *StreamService) createDirs(chunkDir, chunkFileDir, manifestDir, videoDir string) {
	utils.MKDir(chunkDir).Run()
	utils.MKDir(chunkFileDir).Run()
	utils.MKDir(manifestDir).Run()
	utils.MKDir(videoDir).Run()
}

func (ss *StreamService) ServeChunk(filename string) (*schema.OutFile, error) {
	if !ss.chunkHandler.Exists(filename) {
		return nil, ErrChunkNotFound
	}
	return ss.chunkHandler.Retrieve(filename)
}

func (ss *StreamService) ServeEntireVideo(filename string) (*schema.OutFile, error) {
	if !ss.videoHandler.Exists(filename) {
		return nil, ErrVideoNotFound
	}
	return ss.videoHandler.Retrieve(filename)
}

func (ss *StreamService) UploadVideo(file schema.InFile) error {
	return ss.videoHandler.Upload(file)
}

func (ss *StreamService) RemoveVideo(filename string) (err error) {
	if ss.videoHandler.Exists(filename + ".mp4") {
		if err = ss.videoHandler.Remove(filename + ".mp4"); err != nil {
			ss.log.Error(err.Error())
		}
		if err = ss.chunkHandler.Remove(filename); err != nil {
			ss.log.Error(err.Error())
		}
		if err = ss.manifestHandler.Remove(filename + ".m3u8"); err != nil {
			ss.log.Error(err.Error())
		}
	} else {
		return ErrVideoNotFound
	}
	return nil
}
