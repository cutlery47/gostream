package service

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cutlery47/gostream/internal/storage"
	"github.com/cutlery47/gostream/internal/utils"
	"go.uber.org/zap"
)

// service, responsible for all data manipulations
type FileService interface {
	Upload(ctx context.Context, video io.ReadCloser, name string) error
	Remove(ctx context.Context, filename string) error
	Serve(ctx context.Context, filename string) (io.ReadCloser, error)
}

// service for creating and removing manifest file and .ts chunks
type CreatorService interface {
	Create(paths storage.Paths, video io.ReadCloser, name string) (manifest io.ReadCloser, chunks []io.ReadCloser, err error)
	Remove(paths storage.Paths) error
}

// FileService impl
type StreamService struct {
	log            *zap.Logger
	storage        storage.Storage
	creatorService CreatorService
}

func NewStreamService(log *zap.Logger, storage storage.Storage, creatorService CreatorService) *StreamService {
	return &StreamService{
		log:            log,
		storage:        storage,
		creatorService: creatorService,
	}
}

func (ss *StreamService) Upload(ctx context.Context, video *os.File, name string) error {
	// figuring out where to store files locally
	paths := ss.storage.Paths()
	// create necessary directories if don't exist
	createDirs(paths.VidPath, paths.ManPath, paths.ChunkPath, name)

	precisePaths := storage.Paths{
		VidPath:   fmt.Sprintf("%v/%v.mp4", paths.VidPath, name),
		ManPath:   fmt.Sprintf("%v/%v.m3u8", paths.ManPath, name),
		ChunkPath: fmt.Sprintf("%v/%v/", paths.ChunkPath, name),
	}

	// creating all the files locally
	manifest, chunks, err := ss.creatorService.Create(precisePaths, video, name)
	if err != nil {
		return err
	}

	defer ss.creatorService.Remove(precisePaths)

	sVideo, err := storage.FromFD(video, name)
	if err != nil {
		return err
	}

	sManifest, err := storage.FromFD(manifest, manifest.Name())
	if err != nil {
		return err
	}

	var sChunks []storage.File
	for _, chunk := range chunks {
		sChunk, err := storage.FromFD(chunk, chunk.Name())
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

// CreateService impl
type ManifestService struct {
	errLog  *zap.Logger
	infoLog *zap.Logger
}

func NewManifestService(infoLog, errLog *zap.Logger) *ManifestService {
	return &ManifestService{
		errLog:  errLog,
		infoLog: infoLog,
	}
}

func (ms *ManifestService) Create(paths storage.Paths, video *os.File, name string) (*os.File, []*os.File, error) {
	// reading raw .mp4 video file
	videoData, err := io.ReadAll(video)
	if err != nil {
		return nil, nil, err
	}

	// creating .mp4 video file locally
	if err := os.WriteFile(paths.VidPath, videoData, 0664); err != nil {
		return nil, nil, err
	}

	// segmentation + .m3u8 creation
	// results in manifest file and chunks creation
	cmd := utils.SegmentVideoAndCreateManifest(
		paths.VidPath,
		// precise manifest path
		paths.ManPath,
		// chunk file path + template for segmentation
		fmt.Sprintf("%v/%v_%%4d.ts", paths.ChunkPath, name),
	)

	// check if segmentation went smoothely
	out, err := cmd.CombinedOutput()
	if err != nil {
		ms.infoLog.Info(string(out))
		return nil, nil, ErrSegmentationException
	}

	var manifest *os.File
	var chunks []*os.File

	// retrieving manifest data
	manifest, err = os.Open(paths.ManPath)
	if err != nil {
		return nil, nil, err
	}

	// itrating over each chunk in the local directory
	chunkDir, _ := os.ReadDir(paths.ChunkPath)
	// filling up chunk array
	for _, el := range chunkDir {
		// retrieving chunk data
		chunk, err := os.Open(paths.ChunkPath + el.Name())
		if err != nil {
			return nil, nil, err
		}

		chunks = append(chunks, chunk)
	}

	return manifest, chunks, nil
}

func (ms *ManifestService) Remove(paths storage.Paths) error {
	if err := os.Remove(paths.VidPath); err != nil {
		return err
	}

	if err := os.Remove(paths.ManPath); err != nil {
		return err
	}

	if err := os.RemoveAll(paths.ChunkPath); err != nil {
		return err
	}

	return nil
}

func createDirs(vidPath, manPath, chunkPath, objName string) {
	chunkFilePath := fmt.Sprintf("%v/%v", chunkPath, objName)
	utils.MKDir(chunkPath).Run()
	utils.MKDir(chunkFilePath).Run()
	utils.MKDir(manPath).Run()
	utils.MKDir(vidPath).Run()
}
