package service

import (
	"fmt"
	"os"

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
	Create(vidPath, manPath, chunkPath string, file schema.InFile) (manifest *schema.InFile, chunks *[]schema.InFile, err error)
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

func (ms *ManifestService) Create(vidPath, manPath, chunkPath string, file schema.InFile) (manifest *schema.InFile, chunks *[]schema.InFile, err error) {
	// create necessary directories if don't exist
	ms.createDirs(vidPath, manPath, chunkPath, file.Name)

	preciseVidPath := fmt.Sprintf("%v/%v.mp4", vidPath, file.Name)
	preciseManPath := fmt.Sprintf("%v/%v.m3u8", manPath, file.Name)
	preciseChunkDir := fmt.Sprintf("%v/%v/", chunkPath, file.Name)

	// reading raw .mp4 video file
	fileData := make([]byte, file.Size)
	if _, err := file.Raw.Read(fileData); err != nil {
		ms.errLog.Error(err.Error())
	}

	// creating .mp4 video file locally
	if err := os.WriteFile(fmt.Sprintf("%v/%v.mp4", vidPath, file.Name), fileData, 0664); err != nil {
		ms.errLog.Error(err.Error())
	}

	// segmentation + .m3u8 creation
	// results in manifest file and chunks creation
	cmd := utils.SegmentVideoAndCreateManifest(
		preciseVidPath,
		// precise manifest path
		preciseManPath,
		// chunk file path + template for segmentation
		fmt.Sprintf("%v/%v_%%4d.ts", preciseChunkDir, file.Name),
	)

	// check if segmentation went smoothely
	out, err := cmd.CombinedOutput()
	if err != nil {
		ms.infoLog.Info(string(out))
		return nil, nil, ErrSegmentationException
	}

	// creating containers for manifest and chunks
	manifest = &schema.InFile{}
	chunks = &[]schema.InFile{}

	// retrieving manifest data
	manifestFile, _ := os.Open(preciseManPath)
	manifestStat, _ := manifestFile.Stat()

	// filling up manifest container
	manifest.Raw = manifestFile
	manifest.Name = manifestStat.Name()
	manifest.Size = int(manifestStat.Size())

	// itrating over each chunk in the local directory
	chunkFiles, _ := os.ReadDir(preciseChunkDir)
	// filling up chunk array
	for _, chunk := range chunkFiles {
		// retrieving chunk data
		chunkName := chunk.Name()
		chunkFile, _ := os.Open(preciseChunkDir + chunkName)
		chunkStat, _ := chunkFile.Stat()

		// filling up chunk container
		chunkEl := schema.InFile{}
		chunkEl.Name = chunkName
		chunkEl.Raw = chunkFile
		chunkEl.Size = int(chunkStat.Size())

		*chunks = append(*chunks, chunkEl)
	}

	return manifest, chunks, nil
}

func (ss *ManifestService) createDirs(vidPath, manPath, chunkPath, filename string) {
	chunkFilePath := fmt.Sprintf("%v/%v", chunkPath, filename)
	utils.MKDir(chunkPath).Run()
	utils.MKDir(chunkFilePath).Run()
	utils.MKDir(manPath).Run()
	utils.MKDir(vidPath).Run()
}
