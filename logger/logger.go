package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Debug bool
}

func New(cfg Config) zerolog.Logger {
	var level zerolog.Level
	if cfg.Debug {
		level = zerolog.DebugLevel
	} else {
		level = zerolog.InfoLevel
	}
	return log.Level(level)
}

func NewSubLogger(log zerolog.Logger, component string) zerolog.Logger {
	return log.With().Str("component", component).Logger()
}
