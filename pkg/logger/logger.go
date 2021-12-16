package logger

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger is a wrapper for zap.Logger
func NewLogger() *zap.Logger {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// encoderConfig.EnableStackTrace = false                  // sẽ không có stacktracekey

	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	logger := zap.New(core, zap.AddCaller())

	return logger
}
