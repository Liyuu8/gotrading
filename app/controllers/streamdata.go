package controllers

import (
	"gotrading/app/models"
	"gotrading/bitflyer"
	"gotrading/config"
	"log"
)

// StreamIngestionData is ...
func StreamIngestionData() {
	var tickerChannel = make(chan bitflyer.Ticker)
	apiClient := bitflyer.New(config.BitbankConfig.APIKey, config.BitbankConfig.APISecret)
	go apiClient.GetRealTimeTicker(config.BitflyerConfig.ProductCode, tickerChannel)
	for ticker := range tickerChannel {
		log.Printf("action=StreamIngestionData, %v", ticker)
		for _, duration := range config.BitflyerConfig.Durations {
			isCreated := models.CreateCandleWithDuration(ticker, ticker.ProductCode, duration)
			if isCreated && duration == config.BitflyerConfig.TradeDuration {
				// TODO:
			}
		}
	}
}
