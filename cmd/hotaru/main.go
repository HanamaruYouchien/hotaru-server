package main

import (
	"flag"
	"fmt"

	"hotaru.hana.im/server/pkg/config"
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

	config, err := config.ReadConfig(configPath)
	if err != nil {
		fmt.Println("read config failed")
	}
	fmt.Println(config)

	web.Serve()
}
