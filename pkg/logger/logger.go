package logger

import (
	"fmt"
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/xdorro/golang-grpc-base-project/pkg/utils"
)

// NewLogger is a wrapper for zap.Logger
func NewLogger() *zap.Logger {
	encoder := ecszap.NewDefaultEncoderConfig()
	encoder.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder.EncodeDuration = zapcore.MillisDurationEncoder
	encoder.EncodeCaller = ecszap.ShortCallerEncoder

	core := ecszap.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(os.Stdout, getLogWriter()),
		zap.DebugLevel,
	)

	return zap.New(core, zap.AddCaller())
}

// getLogWriter add a lumberjack.Logger to the zap.Logger
func getLogWriter() zapcore.WriteSyncer {
	if err := utils.MakeDir("./logs/"); err != nil {
		panic(fmt.Sprintf("Failed to create logs directory: %s", err))
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/logs.log",
		MaxSize:    10, // MB
		MaxBackups: 5,  // Files
		MaxAge:     30,
		Compress:   false,
	}

	return zapcore.AddSync(lumberJackLogger)
}
