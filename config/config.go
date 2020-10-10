package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

// ConfigList has APIKey and APISecret
type ConfigList struct {
	APIKey    string
	APISecret string
	LogFile   string
}

// Config  is ...
var Config ConfigList

func init() {
	config, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		APIKey:    config.Section("bitbank").Key("api_key").String(),
		APISecret: config.Section("bitbank").Key("api_secret").String(),
		LogFile:   config.Section("gotrading").Key("log_file").String(),
	}
}
