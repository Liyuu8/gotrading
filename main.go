package main

import (
	"./bitbank"
	"./config"
	"./utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	apiClient := bitbank.New(config.Config.APIKey, config.Config.APISecret)
	apiClient.GetAssets()
	apiClient.GetPairs()
	ticker, _ := apiClient.GetTicker("btc_jpy")
	ticker.GetMiddle()
}
