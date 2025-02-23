package logger

import (
	"go.uber.org/zap"
	"log"
)

var sugar *zap.SugaredLogger

// InitLogger - inits the logger
func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(" Zap Logger init error", err)
	}
	sugar = logger.Sugar()
}

// NewLogger - returns the existing logger or creates the new one
func NewLogger() *zap.SugaredLogger {
	if sugar == nil {
		InitLogger()
	}
	return sugar
}
