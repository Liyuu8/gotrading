package config

import (
	"log"
	"os"
	"time"

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

	TradeDuration time.Duration
	Durations     map[string]time.Duration
	DbName        string
	SQLDriver     string
	Port          int
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

	durations := map[string]time.Duration{
		"1s": time.Second,
		"1m": time.Minute,
		"1h": time.Hour,
	}

	BitbankConfig = BitbankConfigList{
		APIKey:    config.Section("bitbank").Key("api_key").String(),
		APISecret: config.Section("bitbank").Key("api_secret").String(),
		LogFile:   config.Section("bitbank").Key("log_file").String(),
		Pair:      config.Section("gotrading").Key("pair").String(),
	}

	BitflyerConfig = BitflyerConfigList{
		APIKey:        config.Section("bitflyer").Key("api_key").String(),
		APISecret:     config.Section("bitflyer").Key("api_secret").String(),
		LogFile:       config.Section("bitflyer").Key("log_file").String(),
		ProductCode:   config.Section("gotrading").Key("product_code").String(),
		Durations:     durations,
		TradeDuration: durations[config.Section("gotrading").Key("trade_duration").String()],
		DbName:        config.Section("db").Key("name").String(),
		SQLDriver:     config.Section("db").Key("driver").String(),
		Port:          config.Section("web").Key("port").MustInt(),
	}
}
