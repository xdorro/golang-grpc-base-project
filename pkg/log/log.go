package log

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
)

func init() {
	encoder := ecszap.NewDefaultEncoderConfig()
	encoder.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder.EncodeDuration = zapcore.MillisDurationEncoder
	encoder.EncodeCaller = ecszap.ShortCallerEncoder

	core := ecszap.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(os.Stdout, getLogWriter()),
		zap.DebugLevel,
	)
	logger = zap.New(
		core,
		zap.AddCaller(),
		// zap.Hooks(func(entry zapcore.Entry) error {
		// fmt.Println("test hooks test hooks")
		// 	return nil
		// }),
	)
}

// NewLogger is a wrapper for zap.Logger
func NewLogger() *zap.Logger {
	return logger
}

// getLogWriter returns a lumberjack.Logger
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/data.log",
		MaxSize:    10, // MB
		MaxBackups: 5,  // Files
		MaxAge:     30,
		Compress:   false,
	}

	return zapcore.AddSync(lumberJackLogger)
}

// Sync is a wrapper for zap.Sync
func Sync() error {
	return logger.Sync()
}
