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
}

type Config struct {
	Log ConfigLog
}

type ConfigLog struct {
	Level  string
	Pretty bool
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
