package service

import (
	"fmt"
	"io"
	"os"

	"github.com/cutlery47/gostream/internal/schema"
	"github.com/cutlery47/gostream/internal/storage"
	"github.com/cutlery47/gostream/internal/utils"
	"go.uber.org/zap"
)

// service, responsible for all data manipulations
type FileService interface {
	Upload(video schema.InVideo) error
	Remove(filename string) error
	Serve(filename string) (*schema.OutFile, error)
}

// service for creating manifest file and .ts chunks
type CreatorService interface {
	Create(vidPath, manPath, chunkPath string, video schema.InVideo) (manifest *schema.InFile, chunks *[]schema.InFile, err error)
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

func (ss *StreamService) Upload(video schema.InVideo) error {
	// figuring out where to store files locally
	paths := ss.storage.Paths()
	manifest, chunks, err := ss.creatorService.Create(paths.VidPath, paths.ManPath, paths.ChunkPath, video)
	if err != nil {
		return err
	}

	return ss.storage.Store(video, *manifest, *chunks)
}

func (ss *StreamService) Remove(filename string) error {
	return ss.storage.Remove(filename)
}

func (ss *StreamService) Serve(filename string) (*schema.OutFile, error) {
	return ss.storage.Get(filename)
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

func (ms *ManifestService) Create(vidPath, manPath, chunkPath string, video schema.InVideo) (manifest *schema.InFile, chunks *[]schema.InFile, err error) {
	// create necessary directories if don't exist
	ms.createDirs(vidPath, manPath, chunkPath, video.Name)

	preciseVidPath := fmt.Sprintf("%v/%v.mp4", vidPath, video.Name)
	preciseManPath := fmt.Sprintf("%v/%v.m3u8", manPath, video.Name)
	preciseChunkDir := fmt.Sprintf("%v/%v/", chunkPath, video.Name)

	// reading raw .mp4 video file
	videoData, err := io.ReadAll(video.File.Raw)
	if err != nil {
		return nil, nil, err
	}

	// creating .mp4 video file locally
	if err := os.WriteFile(fmt.Sprintf("%v/%v.mp4", vidPath, video.Name), videoData, 0664); err != nil {
		return nil, nil, err
	}

	// segmentation + .m3u8 creation
	// results in manifest file and chunks creation
	cmd := utils.SegmentVideoAndCreateManifest(
		preciseVidPath,
		// precise manifest path
		preciseManPath,
		// chunk file path + template for segmentation
		fmt.Sprintf("%v/%v_%%4d.ts", preciseChunkDir, video.Name),
	)

	// check if segmentation went smoothely
	out, err := cmd.CombinedOutput()
	if err != nil {
		ms.infoLog.Info(string(out))
		return nil, nil, ErrSegmentationException
	}

	manifest = &schema.InFile{}
	chunks = &[]schema.InFile{}

	// retrieving manifest data
	manifestFile, _ := os.Open(preciseManPath)
	manifestStat, _ := manifestFile.Stat()

	manifest = &schema.InFile{
		Raw:  manifestFile,
		Name: manifestStat.Name(),
		Size: int(manifestStat.Size()),
	}

	// itrating over each chunk in the local directory
	chunkFiles, _ := os.ReadDir(preciseChunkDir)
	// filling up chunk array
	for _, chunk := range chunkFiles {
		// retrieving chunk data
		chunkName := chunk.Name()
		chunkFile, _ := os.Open(preciseChunkDir + chunkName)
		chunkStat, _ := chunkFile.Stat()

		// filling up chunk container
		chunkEl := schema.InFile{
			Name: chunkName,
			Raw:  chunkFile,
			Size: int(chunkStat.Size()),
		}

		*chunks = append(*chunks, chunkEl)
	}

	return manifest, chunks, nil
}

func (ss *ManifestService) createDirs(vidPath, manPath, chunkPath, objName string) {
	chunkFilePath := fmt.Sprintf("%v/%v", chunkPath, objName)
	utils.MKDir(chunkPath).Run()
	utils.MKDir(chunkFilePath).Run()
	utils.MKDir(manPath).Run()
	utils.MKDir(vidPath).Run()
}
