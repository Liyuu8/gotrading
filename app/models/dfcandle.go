package models

import (
	"time"
)

// DataFrameCandle is ...
type DataFrameCandle struct {
	ProductCode string        `json:"product_code"`
	Duration    time.Duration `json:"duration"`
	Candles     []Candle      `json:"candles"`
}

// Times is ...
func (df *DataFrameCandle) Times() []time.Time {
	s := make([]time.Time, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Time
	}
	return s
}

// Opens is ...
func (df *DataFrameCandle) Opens() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Open
	}
	return s
}

// Closes is ...
func (df *DataFrameCandle) Closes() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Close
	}
	return s
}

// Highs is ...
func (df *DataFrameCandle) Highs() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.High
	}
	return s
}

// Lows is ...
func (df *DataFrameCandle) Lows() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Low
	}
	return s
}

// Volumes is ...
func (df *DataFrameCandle) Volumes() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Volume
	}
	return s
}
