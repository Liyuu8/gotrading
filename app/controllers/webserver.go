package controllers

import (
	"encoding/json"
	"fmt"
	"gotrading/app/models"
	"gotrading/config"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"text/template"
)

func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.ParseFiles("app/views/chart.html"))
	err := templates.ExecuteTemplate(w, "chart.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// JSONError is ...
type JSONError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// APIError is ...
func APIError(w http.ResponseWriter, errMessage string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonError, err := json.Marshal(JSONError{Error: errMessage, Code: code})
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonError)
}

func apiMakeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var apiValidPath = regexp.MustCompile("^/api/candle/$")
		m := apiValidPath.FindStringSubmatch(r.URL.Path)
		if len(m) == 0 {
			APIError(w, "Not found", http.StatusNotFound)
		}
		fn(w, r)
	}
}

func apiCandleHandler(w http.ResponseWriter, r *http.Request) {
	productCode := r.URL.Query().Get("product_code")
	if productCode == "" {
		APIError(w, "No product_code param", http.StatusBadRequest)
		return
	}

	strLimit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(strLimit)
	if strLimit == "" || err != nil || limit < 0 || limit > 1000 {
		limit = 1000
	}

	duration := r.URL.Query().Get("duration")
	if duration == "" {
		duration = "1m"
	}
	durationTime := config.BitflyerConfig.Durations[duration]

	df, _ := models.GetAllCandle(productCode, durationTime, limit)

	sma := r.URL.Query().Get("sma")
	if sma != "" {
		strSmaPeriod1 := r.URL.Query().Get("smaPeriod1")
		strSmaPeriod2 := r.URL.Query().Get("smaPeriod2")
		strSmaPeriod3 := r.URL.Query().Get("smaPeriod3")

		period1, err := strconv.Atoi(strSmaPeriod1)
		if strSmaPeriod1 == "" || err != nil || period1 < 0 {
			period1 = 7
		}
		period2, err := strconv.Atoi(strSmaPeriod2)
		if strSmaPeriod2 == "" || err != nil || period2 < 0 {
			period2 = 14
		}
		period3, err := strconv.Atoi(strSmaPeriod3)
		if strSmaPeriod3 == "" || err != nil || period3 < 0 {
			period3 = 50
		}

		df.AddSma(period1)
		df.AddSma(period2)
		df.AddSma(period3)
	}

	js, err := json.Marshal(df)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// StartWebServer is ...
func StartWebServer() error {
	http.Handle(
		"/resources/",
		http.StripPrefix(
			"/resources/",
			http.FileServer(http.Dir("app/views/resources/")),
		),
	)

	http.HandleFunc("/api/candle/", apiMakeHandler(apiCandleHandler))
	// http://localhost:8080/api/candle/?product_code=BTC_JPY
	// http://localhost:8080/api/candle/?product_code=BTC_JPY&duration=1s
	// http://localhost:8080/api/candle/?product_code=BTC_JPY&duration=1s&limit=1

	http.HandleFunc("/chart/", viewChartHandler)
	// http://localhost:8080/chart/

	return http.ListenAndServe(fmt.Sprintf(":%d", config.BitflyerConfig.Port), nil)
}
