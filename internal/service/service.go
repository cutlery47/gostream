package service

import (
	"github.com/cutlery47/gostream/internal/schema"
	"github.com/cutlery47/gostream/internal/storage"
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
	Create(file schema.InFile) (manifest schema.InFile, chunks []schema.InFile, err error)
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
	manifest, chunks, err := ss.creatorService.Create(video)
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

func (ms *ManifestService) Create(file schema.InFile) (manifest schema.InFile, chunks []schema.InFile, err error) {
	// todo: insert creation logic here
	return schema.InFile{}, []schema.InFile{}, nil
}

// func (ss *StreamService) UploadVideo(file schema.InFile) error {
// 	// 1) receives video file
// 	// 2) uploads raw video, creates manifest with chunks and uplods them as well
// 	if err := ss.videoHandler.Upload(file); err != nil {
// 		return err
// 	}

// 	if err := ss.manifestHandler.CreateAndUpload(file); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (ss *StreamService) ServeManifestOrChunk(filename string) (*schema.OutFile, error) {
// 	return ss.manifestHandler.Retrieve(filename)
// }

// TODO: add chunk checks
// func (ss *StreamService) ServeManifest(filename string) (*schema.OutFile, error) {
// 	// check if we already store the manifest file
// 	if manifest, _ := ss.manifestHandler.Retrieve(filename + ".m3u8"); manifest != nil {
// 		ss.log.Info(fmt.Sprintf("Manifest for %v is already stored. Returning...", filename))
// 		return manifest, nil
// 	}

// 	// check if video file is even stored
// 	if !ss.videoHandler.Exists(fmt.Sprintf("%v.mp4", filename)) {
// 		return nil, ErrVideoNotFound
// 	}

// 	chunkDir := ss.chunkHandler.Path()
// 	chunkFileDir := fmt.Sprintf("%v/%v", chunkDir, filename)
// 	manifestDir := ss.manifestHandler.Path()
// 	videoDir := ss.videoHandler.Path()

// 	ss.createDirs(chunkDir, chunkFileDir, manifestDir, videoDir)
// 	cmd := utils.SegmentVideoAndCreateManifest(
// 		// precise file path
// 		fmt.Sprintf("%v/%v.mp4", videoDir, filename),
// 		// precise manifest path
// 		fmt.Sprintf("%v/%v.m3u8", manifestDir, filename),
// 		// chunk file path + template for segmentation
// 		fmt.Sprintf("%v/%v_%%4d.ts", chunkFileDir, filename),
// 	)

// 	// check if segmentation went smoothely
// 	out, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return nil, ErrSegmentationException
// 	}

// 	ss.log.Error(string(out))

// 	return ss.manifestHandler.Retrieve(filename + ".m3u8")
// }

// func (ss *StreamService) createDirs(chunkDir, chunkFileDir, manifestDir, videoDir string) {
// 	utils.MKDir(chunkDir).Run()
// 	utils.MKDir(chunkFileDir).Run()
// 	utils.MKDir(manifestDir).Run()
// 	utils.MKDir(videoDir).Run()
// }

// func (ss *StreamService) ServeChunk(filename string) (*schema.OutFile, error) {
// 	if !ss.chunkHandler.Exists(filename) {
// 		return nil, ErrChunkNotFound
// 	}
// 	return ss.chunkHandler.Retrieve(filename)
// }

// func (ss *StreamService) ServeEntireVideo(filename string) (*schema.OutFile, error) {
// 	if !ss.videoHandler.Exists(filename) {
// 		return nil, ErrVideoNotFound
// 	}
// 	return ss.videoHandler.Retrieve(filename)
// }

// func (ss *StreamService) UploadVideo(file schema.InFile) error {
// 	return ss.videoHandler.Upload(file)
// }

// func (ss *StreamService) RemoveVideo(filename string) (err error) {
// 	if ss.videoHandler.Exists(filename + ".mp4") {
// 		if err = ss.videoHandler.Remove(filename + ".mp4"); err != nil {
// 			ss.log.Error(err.Error())
// 		}
// 		if err = ss.chunkHandler.Remove(filename); err != nil {
// 			ss.log.Error(err.Error())
// 		}
// 		if err = ss.manifestHandler.Remove(filename + ".m3u8"); err != nil {
// 			ss.log.Error(err.Error())
// 		}
// 	} else {
// 		return ErrVideoNotFound
// 	}
// 	return nil
// }
