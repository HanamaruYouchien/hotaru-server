package main

import (
	"os"

	"github.com/goccy/go-yaml"
)

var defaultConfig = Config{
	Log: ConfigLog{
		Level:  "info",
		Pretty: true,
	},
	Database: ConfigDatabase{
		Type: "sqlite",
		Url:  ":memory:",
	},
}

type Config struct {
	Log      ConfigLog
	Database ConfigDatabase
}

type ConfigLog struct {
	Level  string
	Pretty bool
}

type ConfigDatabase struct {
	Type string
	Url  string
}

func ReadConfig(configPath string) (config Config, err error) {
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return defaultConfig, err
	}

	if err = yaml.Unmarshal(configFile, &config); err != nil {
		return defaultConfig, err
	}
	return
}
