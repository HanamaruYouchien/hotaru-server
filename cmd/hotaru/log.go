package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetLogger(configLog ConfigLog) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	level, err := zerolog.ParseLevel(configLog.Level)
	zerolog.SetGlobalLevel(level)

	if configLog.Pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	if err != nil {
		log.Warn().Err(err).Msg("parse level failed")
	}
}
