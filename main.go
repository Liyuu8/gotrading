package main

import (
	"fmt"
	"gotrading/app/controllers"

	"gotrading/app/models"
	"gotrading/config"
	"gotrading/utils"
)

func main() {
	utils.LoggingSettings(config.BitflyerConfig.LogFile)
	fmt.Println(models.DBConnection)
	controllers.StreamIngestionData()
}
