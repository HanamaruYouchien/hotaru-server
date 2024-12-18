package main

import (
	"flag"
	"fmt"

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
	server.Serve()
}

func ptr[T any](x T) *T {
	return &x
}
