package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gotrading/config"

	_ "github.com/mattn/go-sqlite3"
)

const (
	tableNameSignalEvents = "signal_events"
)

// DBConnection is ...
var DBConnection *sql.DB

// GetCandleTableName is ...
func GetCandleTableName(productCode string, duration time.Duration) string {
	return fmt.Sprintf("%s_%s", productCode, duration)
}
func init() {
	var err error
	DBConnection, err = sql.Open(config.BitflyerConfig.SQLDriver, config.BitflyerConfig.DbName)
	if err != nil {
		log.Fatalln(err)
	}
	cmd := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			time DATETIME PRIMARY KEY NOT NULL,
			product_code STRING,
			side STRING,
			price FLOAT,
			size FLOAT
		)
	`, tableNameSignalEvents)
	DBConnection.Exec(cmd)

	for _, duration := range config.BitflyerConfig.Durations {
		tableName := GetCandleTableName(config.BitflyerConfig.ProductCode, duration)
		c := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				time DATETIME PRIMARY KEY NOT NULL,
				open FLOAT,
				close FLOAT,
				high FLOAT,
				low FLOAT,
				volume FLOAT
			)
		`, tableName)
		DBConnection.Exec(c)
	}
}
