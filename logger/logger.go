package logger

import (
	"io"
	"log"
	"os"

	"log/slog"

	"github.com/gaydin/journey/configuration"
)

func New(config *configuration.Configuration) (*slog.Logger, func() error) {
	var (
		writer           io.Writer
		logFileCloseFunc = func() error { return nil }
	)

	if config.Logger.FileName != "" {
		logFile, err := os.OpenFile(config.Logger.FileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Error: Couldn't open log file: " + err.Error())
		}

		logFileCloseFunc = logFile.Close
		writer = logFile
	} else {
		writer = os.Stdout
	}

	var handler slog.Handler
	if config.Logger.JSON {
		handler = slog.NewJSONHandler(writer, nil)
	} else {
		handler = slog.NewTextHandler(writer, nil)
	}

	return slog.New(handler), logFileCloseFunc
}
