package log

import (
	"os"
	"strconv"
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

	// Caller Marshal Function
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	Logger = zerolog.New(mw).With().
		Timestamp().
		Caller().
		Logger()
}

// getLogWriter returns a lumberjack.logger
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
