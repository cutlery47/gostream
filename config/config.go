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
	Distr       distrStorageConfig
}

type localStorageConfig struct {
	ManifestPath string
	ChunkPath    string
	VideoPath    string
}

type distrStorageConfig struct {
	S3Config S3Config
	DBConfig DBConfig
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

type S3Config struct {
}

type segmentConfig struct {
	Time int
}

func New() (*Config, error) {
	godotenv.Load(".env")

	logConfig := loggerConfig{
		RequestLogsPath: os.Getenv("REQUEST_LOGS_PATH"),
		ErrorLogsPath:   os.Getenv("ERROR_LOGS_PATH"),
		InfoLogsPath:    os.Getenv("INFO_LOGS_PATH"),
	}

	lsConfig := localStorageConfig{
		ManifestPath: os.Getenv("MANIFEST_PATH"),
		ChunkPath:    os.Getenv("CHUNK_PATH"),
		VideoPath:    os.Getenv("VIDEO_PATH"),
	}

	s3Config := S3Config{}

	dbConfig := DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
	}

	dsConfig := distrStorageConfig{
		S3Config: s3Config,
		DBConfig: dbConfig,
	}

	sConfig := storageConfig{
		StorageType: os.Getenv("STORAGE_TYPE"),
		Local:       lsConfig,
		Distr:       dsConfig,
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
