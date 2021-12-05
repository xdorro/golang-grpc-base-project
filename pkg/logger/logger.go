package logger

import (
	"log"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger is a wrapper for zap.Logger
func NewLogger() *zap.Logger {
	// cfg := zap.Config{
	// 	Encoding: "console",                           // encode kiểu json hoặc console
	// 	Level:    zap.NewAtomicLevelAt(zap.InfoLevel), // chọn InfoLevel có thể log ở cả 3 level

	// logger, err := cfg.Build() // Build ra Logger

	var cfg zap.Config
	if viper.GetBool("APP_DEBUG") {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // Lấy dòng code bắt đầu log
	cfg.EncoderConfig.EncodeLevel = CustomLevelEncoder          // Format cách hiển thị level log
	cfg.EncoderConfig.EncodeTime = SyslogTimeEncoder            // Format hiển thị thời điểm log

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("create logger: %v", err)
	}

	return logger
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}
