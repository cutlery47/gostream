package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// shamelessly copied from https://betterstack.com/community/guides/logging/go/zap/#examining-zap-s-logging-api
// (idk what halfa dis does bru)
func New(filepath string, isErr bool) *zap.Logger {
	fd, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("logger.New (path: %v): %v", filepath, err)
		return nil
	}

	stdout := zapcore.AddSync(os.Stdout)
	file := zapcore.AddSync(fd)

	var level zap.AtomicLevel
	if isErr {
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	} else {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	prodConfig := zap.NewProductionEncoderConfig()
	prodConfig.TimeKey = "timestamp"
	prodConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	devConfig := zap.NewDevelopmentEncoderConfig()
	devConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(devConfig)
	fileEncoder := zapcore.NewJSONEncoder(prodConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	return zap.New(core)
}
