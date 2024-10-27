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
	AppLogsPath string
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
	Host        string
	Port        string
	AdmPort     string
	VidBucket   string
	ManBucket   string
	ChunkBucket string
	User        string
	Password    string
}

type segmentConfig struct {
	Time int
}

func New() (*Config, error) {
	godotenv.Load("yours.env")

	logConfig := loggerConfig{
		AppLogsPath: os.Getenv("APP_LOGS_PATH"),
	}

	lsConfig := localStorageConfig{
		ManifestPath: os.Getenv("MANIFEST_PATH"),
		ChunkPath:    os.Getenv("CHUNK_PATH"),
		VideoPath:    os.Getenv("VIDEO_PATH"),
	}

	s3Config := S3Config{
		Host:        os.Getenv("MINIO_HOST"),
		Port:        os.Getenv("MINIO_PORT"),
		AdmPort:     os.Getenv("MINIO_ADMIN_PORT"),
		VidBucket:   os.Getenv("MINIO_VID_BUCKET"),
		ChunkBucket: os.Getenv("MINIO_CHUNK_BUCKET"),
		ManBucket:   os.Getenv("MINIO_MAN_BUCKET"),
	}

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
