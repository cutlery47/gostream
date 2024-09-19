package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Log     loggerConfig
	Storage storageConfig
	Segment segmentConfig
}

type loggerConfig struct {
	RequestLogsPath string
	ErrorLogsPath   string
	InfoLogsPath    string
}

type storageConfig struct {
	StorageType string
	Local       localStorageConfig
}

type localStorageConfig struct {
	ManifestPath string
	ChunkPath    string
	VideoPath    string
}

type segmentConfig struct {
	Time int
}

func New() (*Config, error) {
	godotenv.Load(".env")

	lsConfig := localStorageConfig{
		ManifestPath: os.Getenv("MANIFEST_PATH"),
		ChunkPath:    os.Getenv("CHUNK_PATH"),
		VideoPath:    os.Getenv("VIDEO_PATH"),
	}

	sConfig := storageConfig{
		StorageType: os.Getenv("STORAGE_TYPE"),
		Local:       lsConfig,
	}

	logConfig := loggerConfig{
		RequestLogsPath: os.Getenv("REQUEST_LOGS_PATH"),
		ErrorLogsPath:   os.Getenv("ERROR_LOGS_PATH"),
		InfoLogsPath:    os.Getenv("INFO_LOGS_PATH"),
	}

	segtime, err := strconv.Atoi(os.Getenv("SEGMENT_TIME"))
	if err != nil {
		return nil, err
	}

	segConfig := segmentConfig{
		Time: segtime,
	}

	return &Config{
		Log:     logConfig,
		Storage: sConfig,
		Segment: segConfig,
	}, nil
}
