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
	encodeCfg := zap.NewProductionEncoderConfig()
	encodeCfg.EncodeCaller = zapcore.ShortCallerEncoder // Lấy dòng code bắt đầu log
	encodeCfg.EncodeLevel = CustomLevelEncoder          // Format cách hiển thị level log
	encodeCfg.EncodeTime = SyslogTimeEncoder            // Format hiển thị thời điểm log
	encodeCfg.StacktraceKey = ""                        // sẽ không có stacktracekey

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel) // chọn InfoLevel có thể log ở cả 3 level
	cfg.EncoderConfig = encodeCfg

	if viper.GetBool("APP_DEBUG") {
		cfg.Encoding = "console"
	}

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
