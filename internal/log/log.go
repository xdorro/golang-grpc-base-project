package log

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack/v3"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

// init initializes the logger
func init() {
	// UNIX Time is faster and smaller than most timestamps
	consoleWriter := &zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	// Multi Writer
	mw := zerolog.MultiLevelWriter(consoleWriter, getLogWriter())
	Logger = zerolog.New(mw).With().
		Timestamp().
		Caller().
		Logger()
}

// getLogWriter returns a lumberjack.Logger
func getLogWriter() *lumberjack.Roller {
	options := &lumberjack.Options{
		MaxBackups: 5,  // Files
		MaxAge:     30, // 30 days
		Compress:   false,
	}

	roller, err := lumberjack.NewRoller(
		"./logs/data.log",
		500*1024*1024, // 500 MB
		options,
	)

	if err != nil {
		panic(err)
	}

	return roller
}
