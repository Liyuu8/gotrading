package controllers

import (
	"fmt"
	"gotrading/app/models"
	"gotrading/config"
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseFiles("app/views/google.html"))

func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	limit := 100
	duration := "1m"
	durationTime := config.BitflyerConfig.Durations[duration]
	df, _ := models.GetAllCandle(config.BitflyerConfig.ProductCode, durationTime, limit)

	err := templates.ExecuteTemplate(w, "google.html", df.Candles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// StartWebServer is ...
func StartWebServer() error {
	http.HandleFunc("/chart/", viewChartHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.BitflyerConfig.Port), nil)
}
