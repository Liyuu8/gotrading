package main

import (
	"gotrading/app/controllers"

	"gotrading/config"
	"gotrading/utils"
)

func main() {
	utils.LoggingSettings(config.BitflyerConfig.LogFile)
	controllers.StreamIngestionData()
	controllers.StartWebServer()
}
