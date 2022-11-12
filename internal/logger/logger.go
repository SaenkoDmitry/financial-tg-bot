package logger

import (
	"go.uber.org/zap"
	"log"
)

var logger *zap.Logger

func init() {
	localLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("logger init", err)
	}

	logger = localLogger
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
