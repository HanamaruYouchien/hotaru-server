package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"hotaru.hana.im/server/pkg/storage"
	"hotaru.hana.im/server/pkg/web"
)

var Version = "dev"

var DefaultConfigPath string

func init() {
	if Version == "dev" {
		DefaultConfigPath = "../../config/hotaru-server.yml"
	} else {
		DefaultConfigPath = "/etc/hotaru-server.yml"
	}
}

func main() {
	fmt.Println("Hotaru server - " + Version)

	var configPath string
	var firstRun bool
	flag.StringVar(&configPath, "config", DefaultConfigPath, "config path")
	flag.BoolVar(&firstRun, "init", false, "Initialize database")
	flag.Parse()

	// parse config
	config, err := ReadConfig(configPath)
	SetLogger(config.Log)
	if err != nil {
		log.Error().Err(err).Msg("failed to read config, use default config instead")
	}
	log.Debug().Any("config", config).Send()

	// connect database
	database, err := storage.NewStorage(config.Database.Type, config.Database.Url, ptr(log.With().Str("comp", "storage").Logger()))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}
	if firstRun {
		if err := database.Init(); err != nil {
			log.Fatal().Err(err).Msg("failed to init database")
		}
		log.Info().Msg("database init ok")
	}

	// serve http
	server, err := web.NewServer(database, ptr(log.With().Str("comp", "web").Logger()))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create web server")
	}

	go func() {
		if err := server.Serve(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error().Err(err).Msg("web server error")
			}
		}
	}()

	// shutdown gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	if err := server.Shutdown(); err != nil {
		log.Error().Err(err).Msg("shutdown web server error")
	}
	log.Info().Msg("web server shutdown gracefully")
}

func ptr[T any](x T) *T {
	return &x
}
