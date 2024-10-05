package service

import (
	"io"

	"github.com/cutlery47/gostream/internal/storage"
	"go.uber.org/zap"
)

type Retriever interface {
	Retrieve(filename string) (io.Reader, error)
	Exists(filename string) bool
	Path() string
}

type Uploader interface {
	Upload(file io.Reader, filename string) error
}

type Remover interface {
	Remove(filename string) error
}

type RemoveRetriever interface {
	Retriever
	Remover
}

type UploadRemoveRetriever interface {
	Retriever
	Uploader
	Remover
}

type manifestHandler struct {
	log     *zap.Logger
	storage storage.Storage
}

func NewManifestHandler(log *zap.Logger, storage storage.Storage) *manifestHandler {
	return &manifestHandler{
		log:     log,
		storage: storage,
	}
}

func (mh *manifestHandler) Retrieve(filename string) (io.Reader, error) {
	if !mh.storage.Exists(filename) {
		return nil, ErrManifestNotFound
	}
	return mh.storage.Get(filename)
}

func (mh *manifestHandler) Exists(filename string) bool {
	return mh.storage.Exists(filename)
}

func (mh *manifestHandler) Path() string {
	return mh.storage.Path()
}

func (mh *manifestHandler) Remove(filename string) error {
	if !mh.storage.Exists(filename) {
		return ErrManifestNotFound
	}
	return mh.storage.Remove(filename)
}

type chunkHandler struct {
	log     *zap.Logger
	storage storage.Storage
}

func NewChunkHandler(infoLog *zap.Logger, storage storage.Storage) *chunkHandler {
	return &chunkHandler{
		log:     infoLog,
		storage: storage,
	}
}

func (ch *chunkHandler) Retrieve(filename string) (io.Reader, error) {
	if !ch.storage.Exists(filename) {
		return nil, ErrChunkNotFound
	}
	return ch.storage.Get(filename)
}

func (ch *chunkHandler) Exists(filename string) bool {
	return ch.storage.Exists(filename)
}

func (ch *chunkHandler) Path() string {
	return ch.storage.Path()
}

// removes all contents of a chunk directory with a filename
func (ch *chunkHandler) Remove(filename string) error {
	return ch.storage.Remove(filename)
}

type videoHandler struct {
	log     *zap.Logger
	storage storage.Storage
}

func NewVideoHandler(infoLog *zap.Logger, storage storage.Storage) *videoHandler {
	return &videoHandler{
		log:     infoLog,
		storage: storage,
	}
}

func (vh *videoHandler) Retrieve(filename string) (io.Reader, error) {
	if !vh.storage.Exists(filename) {
		return nil, ErrVideoNotFound
	}
	return vh.storage.Get(filename)
}

func (vh *videoHandler) Exists(filename string) bool {
	return vh.storage.Exists(filename)
}

func (vh *videoHandler) Path() string {
	return vh.storage.Path()
}

func (vh *videoHandler) Upload(file io.Reader, filename string) error {
	return vh.storage.Store(file, filename)
}

func (vh *videoHandler) Remove(filename string) error {
	if !vh.storage.Exists(filename) {
		return ErrVideoNotFound
	}
	return vh.storage.Remove(filename)
}
