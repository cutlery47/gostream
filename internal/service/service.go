package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cutlery47/gostream/config"
	"github.com/cutlery47/gostream/internal/storage"
	"github.com/cutlery47/gostream/internal/utils"
	"go.uber.org/zap"
)

// service, responsible for all data manipulations
type Service interface {
	Upload(ctx context.Context, videoReader io.ReadCloser, videoName string) error
	Remove(ctx context.Context, filename string) error
	Serve(ctx context.Context, filename string) (io.ReadCloser, error)
}

type StreamService struct {
	storage storage.Storage

	cfg config.LocalConfig
	log *zap.Logger
}

func NewStreamService(log *zap.Logger, cfg config.LocalConfig, storage storage.Storage) *StreamService {
	return &StreamService{
		storage: storage,

		cfg: cfg,
		log: log,
	}
}

func (ss *StreamService) Upload(ctx context.Context, videoReader io.ReadCloser, videoName string) error {
	// create necessary directories if don't exist
	createDirs(ss.cfg.VideoPath, ss.cfg.ManifestPath, ss.cfg.ChunkPath, videoName)

	videoPath := fmt.Sprintf("%v/%v.mp4", ss.cfg.VideoPath, videoName)
	video, err := createVideo(videoReader, videoPath)
	if err != nil {
		return err
	}

	// creating all the files locally
	manifestPath := fmt.Sprintf("%v/%v.m3u8", ss.cfg.ManifestPath, videoName)
	chunkPath := fmt.Sprintf("%v/%v/", ss.cfg.ChunkPath, videoName)
	manifest, chunks, err := createManifestAndChunks(ss.log, manifestPath, chunkPath, videoPath, videoName)
	if err != nil {
		return err
	}

	// values to be filled and passed to the storage
	var sVideo *storage.File
	var sManifest *storage.File
	var sChunks []storage.File

	nameFromPath := func(path string) string {
		pathSlice := strings.Split(path, "/")
		return pathSlice[len(pathSlice)-1]
	}

	if sVideo, err = storage.FromFD(video, videoName); err != nil {
		return err
	}

	if sManifest, err = storage.FromFD(manifest, nameFromPath(manifest.Name())); err != nil {
		return err
	}

	for _, chunk := range chunks {
		sChunk, err := storage.FromFD(chunk, nameFromPath(chunk.Name()))
		if err != nil {
			return err
		}

		sChunks = append(sChunks, *sChunk)
	}

	return ss.storage.Store(ctx, *sVideo, *sManifest, sChunks)
}

func (ss *StreamService) Remove(ctx context.Context, filename string) error {
	return ss.storage.Remove(ctx, filename)
}

func (ss *StreamService) Serve(ctx context.Context, filename string) (io.ReadCloser, error) {
	return ss.storage.Get(ctx, filename)
}

func createVideo(videoReader io.ReadCloser, videoPath string) (*os.File, error) {
	// reading raw .mp4 video file
	videoData, err := io.ReadAll(videoReader)
	if err != nil {
		return nil, err
	}

	video, err := os.OpenFile(videoPath, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		return nil, err
	}

	if _, err := video.Write(videoData); err != nil {
		return nil, err
	}

	return video, nil
}

func createManifestAndChunks(infoLog *zap.Logger, manifestPath, chunkPath, videoPath, videoName string) (*os.File, []*os.File, error) {
	// segmentation + .m3u8 creation
	// results in manifest file and chunks creation
	cmd := utils.SegmentVideoAndCreateManifest(
		videoPath,
		// precise manifest path
		manifestPath,
		// chunk file path + template for segmentation
		fmt.Sprintf("%v/%v_%%4d.ts", chunkPath, videoName),
	)

	// check if segmentation went smoothely
	out, err := cmd.CombinedOutput()
	if err != nil {
		infoLog.Info(string(out))
		return nil, nil, ErrSegmentationException
	}

	var manifest *os.File
	var chunks []*os.File

	// retrieving manifest data
	manifest, err = os.Open(manifestPath)
	if err != nil {
		return nil, nil, err
	}

	// itrating over each chunk in the local directory
	chunkDir, _ := os.ReadDir(chunkPath)
	// filling up chunk array
	for _, el := range chunkDir {
		// retrieving chunk data
		chunk, err := os.Open(chunkPath + el.Name())
		if err != nil {
			return nil, nil, err
		}

		chunks = append(chunks, chunk)
	}

	return manifest, chunks, nil
}

func createDirs(vidPath, manPath, chunkPath, objName string) {
	chunkFilePath := fmt.Sprintf("%v/%v", chunkPath, objName)
	utils.MKDir(chunkPath).Run()
	utils.MKDir(chunkFilePath).Run()
	utils.MKDir(manPath).Run()
	utils.MKDir(vidPath).Run()
}
