package service

import (
	"fmt"

	"github.com/cutlery47/gostream/internal/schema"
	"github.com/cutlery47/gostream/internal/storage"
	"github.com/cutlery47/gostream/internal/utils"
	"go.uber.org/zap"
)

// service, responsible for all data manipulations
type FileService interface {
	Upload(video schema.InFile) error
	Remove(filename string) error
	Serve(filename string) (*schema.OutFile, error)
}

// service for creating manifest file and .ts chunks
type CreatorService interface {
	Create(vidPath, manPath, chunkPath, filename string) (manifest schema.InFile, chunks []schema.InFile, err error)
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

func (ss *StreamService) Upload(video schema.InFile) error {
	// figuring out where to store files locally
	paths := ss.storage.Paths()
	manifest, chunks, err := ss.creatorService.Create(paths.VidPath, paths.ManPath, paths.ChunkPath, video.Name)
	if err != nil {
		return err
	}

	return ss.storage.Store(video, manifest, chunks)
}

func (ss *StreamService) Remove(filename string) error {
	return ss.storage.Remove(filename)
}

func (ss *StreamService) Serve(filename string) (*schema.OutFile, error) {
	return ss.storage.Get(filename)
}

// CreateService impl
type ManifestService struct {
	log *zap.Logger
}

func NewManifestService(log *zap.Logger) *ManifestService {
	return &ManifestService{
		log: log,
	}
}

func (ms *ManifestService) Create(vidPath, manPath, chunkPath, filename string) (manifest schema.InFile, chunks []schema.InFile, err error) {
	// create necessary directories if don't exist
	ms.createDirs(vidPath, manPath, chunkPath, filename)

	// segmentation + .m3u8 creation
	cmd := utils.SegmentVideoAndCreateManifest(
		// precise file path
		fmt.Sprintf("%v/%v.mp4", vidPath, filename),
		// precise manifest path
		fmt.Sprintf("%v/%v.m3u8", manPath, filename),
		// chunk file path + template for segmentation
		fmt.Sprintf("%v/%v_%%4d.ts", chunkPath, filename),
	)

	// check if segmentation went smoothely
	out, err := cmd.CombinedOutput()
	if err != nil {
		return schema.InFile{}, []schema.InFile{}, ErrSegmentationException
	}

	ms.log.Info(string(out))

	return schema.InFile{}, []schema.InFile{}, nil
}

func (ss *ManifestService) createDirs(vidPath, manPath, chunkPath, filename string) {
	chunkFilePath := fmt.Sprintf("%v/%v", chunkPath, filename)
	utils.MKDir(chunkPath).Run()
	utils.MKDir(chunkFilePath).Run()
	utils.MKDir(manPath).Run()
	utils.MKDir(vidPath).Run()
}
