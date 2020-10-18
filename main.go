package main

import (
	"fmt"

	"./bitflyer"
	"./config"
	"./utils"
)

func main() {
	utils.LoggingSettings(config.BitflyerConfig.LogFile)
	apiClient := bitflyer.New(config.BitflyerConfig.APIKey, config.BitflyerConfig.APISecret)

	// apiClient.GetBalance()
	// ticker, _ := apiClient.GetTicker("BTC_JPY")
	// ticker.GetMidPrice()
	// ticker.DateTime()
	// ticker.TruncateDateTime(time.Hour)

	// tickerChannel := make(chan bitflyer.Ticker)
	// go apiClient.GetRealTimeTicker(config.BitflyerConfig.ProductCode, tickerChannel)
	// for ticker := range tickerChannel {
	// 	fmt.Println(ticker)
	// 	fmt.Println(ticker.GetMidPrice())
	// 	fmt.Println(ticker.DateTime())
	// 	fmt.Println(ticker.TruncateDateTime(time.Second))
	// 	fmt.Println(ticker.TruncateDateTime(time.Minute))
	// 	fmt.Println(ticker.TruncateDateTime(time.Hour))
	// }

	res, _ := apiClient.SendOrder(&bitflyer.Order{
		ProductCode:     config.BitflyerConfig.ProductCode,
		ChildOrderType:  "MARKET",
		Side:            "SELL", // "BUY" or "SELL"
		Size:            0.01,
		MinuteToExpires: 1,
		TimeInForce:     "GTC",
	})
	fmt.Println(res.ChildOrderAcceptanceID)

	// res, _ := apiClient.GetOrders(map[string]string{
	// 	"product_code":              config.BitflyerConfig.ProductCode,
	// 	"child_order_acceptance_id": "JRF20201018-124254-204542",
	// })
	// fmt.Println(res)
}
