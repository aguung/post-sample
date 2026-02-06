package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(env string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if env == "dev" || env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		// JSON output for production
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}
}

func GetLogger() zerolog.Logger {
	return log.Logger
}
