package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Log     LoggerConfig
	Storage StorageConfig
	Flag    FlagConfig
}

type LoggerConfig struct {
	AppLogsPath string `env:"APP_LOGS_PATH"`
}

type StorageConfig struct {
	Local LocalConfig
	Distr DistrConfig
}

type LocalConfig struct {
	ManifestPath string `env:"MANIFEST_PATH"`
	ChunkPath    string `env:"CHUNK_PATH"`
	VideoPath    string `env:"VIDEO_PATH"`
}

type DistrConfig struct {
	S3Config S3Config
	DBConfig DBConfig
}

type DBConfig struct {
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	DBName   string `env:"POSTGRES_NAME"`
	SSLMode  string `env:"POSTGRES_SSL"`
}

type S3Config struct {
	Host        string `env:"MINIO_HOST"`
	Port        string `env:"MINIO_PORT"`
	AdmPort     string `env:"MINIO_ADMIN_PORT"`
	VidBucket   string `env:"MINIO_VID_BUCKET"`
	ManBucket   string `env:"MINIO_CHUNK_BUCKET"`
	ChunkBucket string `env:"MINIO_MAN_BUCKET"`
}

type FlagConfig struct {
	Time string `env:"SEGMENT_TIME"`
	Type string `env:"STORAGE_TYPE"`
}

func New() (cfg *Config, err error) {
	godotenv.Load("yours.env")

	var s3Conf S3Config
	var dbConf DBConfig
	var locConf LocalConfig
	var logConf LoggerConfig
	var flgConf FlagConfig

	confs := []interface{}{&s3Conf, &dbConf, &locConf, &logConf, &flgConf}
	for _, conf := range confs {
		if err = cleanenv.ReadEnv(conf); err != nil {
			return nil, err
		}
	}

	cfg = &Config{
		Log: logConf,
		Storage: StorageConfig{
			Local: locConf,
			Distr: DistrConfig{
				DBConfig: dbConf,
				S3Config: s3Conf,
			},
		},
		Flag: flgConf,
	}

	return cfg, nil
}
