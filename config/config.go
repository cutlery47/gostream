package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Log loggerConfig
}

type loggerConfig struct {
	RequestLogsPath string
	ErrorLogsPath   string
	InfoLogsPath    string
}

func New() *Config {
	godotenv.Load(".env")

	log := loggerConfig{
		RequestLogsPath: os.Getenv("REQUEST_LOGS_PATH"),
		ErrorLogsPath:   os.Getenv("ERROR_LOGS_PATH"),
		InfoLogsPath:    os.Getenv("INFO_LOGS_PATH"),
	}

	return &Config{
		Log: log,
	}
}
