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

	config, err := ReadConfig(configPath)
	SetLogger(config.Log)
	if err != nil {
		log.Error().Err(err).Msg("read config failed, use default config instead")
	}
	log.Debug().Any("config", config).Send()

	server := web.NewServer(ptr(log.With().Str("comp", "web").Logger()))
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
