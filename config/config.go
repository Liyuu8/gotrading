package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

// BitbankConfigList has APIKey and APISecret
type BitbankConfigList struct {
	APIKey    string
	APISecret string
	LogFile   string
	Pair      string
}

// BitflyerConfigList  is ...
type BitflyerConfigList struct {
	APIKey      string
	APISecret   string
	LogFile     string
	ProductCode string
}

// BitbankConfig  is ...
var BitbankConfig BitbankConfigList

// BitflyerConfig  is ...
var BitflyerConfig BitflyerConfigList

func init() {
	config, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	BitbankConfig = BitbankConfigList{
		APIKey:    config.Section("bitbank").Key("api_key").String(),
		APISecret: config.Section("bitbank").Key("api_secret").String(),
		LogFile:   config.Section("bitbank").Key("log_file").String(),
		Pair:      config.Section("gotrading").Key("pair").String(),
	}

	BitflyerConfig = BitflyerConfigList{
		APIKey:      config.Section("bitflyer").Key("api_key").String(),
		APISecret:   config.Section("bitflyer").Key("api_secret").String(),
		LogFile:     config.Section("bitflyer").Key("log_file").String(),
		ProductCode: config.Section("gotrading").Key("product_code_eth").String(),
	}
}
