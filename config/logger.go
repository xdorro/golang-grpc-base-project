package config

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the logger
func InitLogger() {
	// UNIX Time is faster and smaller than most timestamps
	consoleWriter := &zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	// Format Message
	consoleWriter.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("msg=\"%s\"", i)
	}

	// Multi Writer
	mw := zerolog.MultiLevelWriter(consoleWriter, getLogWriter())
	log.Logger = zerolog.New(mw).With().
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
