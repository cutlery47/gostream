package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Log     loggerConfig
	Storage storageConfig
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
}

func New() *Config {
	godotenv.Load(".env")

	lsConfig := localStorageConfig{
		ManifestPath: os.Getenv("MANIFEST_PATH"),
		ChunkPath:    os.Getenv("CHUNK_PATH"),
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

	return &Config{
		Log:     logConfig,
		Storage: sConfig,
	}
}
