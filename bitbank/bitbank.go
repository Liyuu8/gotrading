package bitbank

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const publicEndpoint = "https://public.bitbank.cc/"
const restEndpoint = "https://api.bitbank.cc/v1/"

// APIClient is ...
type APIClient struct {
	key        string
	secret     string
	httpClient *http.Client
}

// New is ...
func New(key, secret string) *APIClient {
	apiClient := &APIClient{key, secret, &http.Client{}}
	return apiClient
}

func (api APIClient) header(method, endpoint string, body []byte) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	var message string
	if method == "GET" {
		message = timestamp + endpoint + string(body)
	} else {
		message = timestamp + string(body)
	}
	mac := hmac.New(sha256.New, []byte(api.secret))
	mac.Write([]byte(message))
	sign := hex.EncodeToString(mac.Sum(nil))

	return map[string]string{
		"ACCESS-KEY":       api.key,
		"ACCESS-NONCE":     timestamp,
		"ACCESS-SIGNATURE": sign,
	}
}

func (api *APIClient) doRequest(baseURLStr, method, urlPath string, query map[string]string, data []byte) (body []byte, err error) {
	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		return
	}
	apiURL, err := url.Parse(urlPath)
	if err != nil {
		return
	}
	endpoint := baseURL.ResolveReference(apiURL).String()
	log.Printf("action=APIClient.doRequest, endpoint=%s", endpoint)
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	for key, value := range api.header(method, req.URL.RequestURI(), data) {
		req.Header.Add(key, value)
	}
	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Assets is ...
type Assets struct {
	Asset           string `json:"asset"`            // アセット名
	FreeAmount      string `json:"free_amount"`      // 利用可能な量
	AmountPrecision int    `json:"amount_precision"` // 精度
	OnhandAmount    string `json:"onhand_amount"`    // 保有量
	LockedAmount    string `json:"locked_amount"`    // ロックされている量
	WithdrawalFee   string `json:"withdrawal_fee"`   // 引き出し手数料
	StopDeposit     bool   `json:"stop_deposit"`     // 入金ステータス
	StopWithdrawal  bool   `json:"stop_withdrawal"`  // 出金ステータス
}

// GetAssets is ...
func (api *APIClient) GetAssets() ([]Assets, error) {
	url := "user/assets"
	resp, err := api.doRequest(restEndpoint, "GET", url, map[string]string{}, nil)
	log.Printf("action=Assets.GetAssets, url=%s, resp=%s", url, string(resp))
	if err != nil {
		log.Printf("action=Assets.GetAssets, err=%s", err.Error())
		return nil, err
	}
	var assets []Assets
	err = json.Unmarshal(resp, &assets)
	if err != nil {
		log.Printf("action=Assets.GetAssets, err=%s", err.Error())
		return nil, err
	}
	return assets, nil
}

// Pairs is ...
type Pairs struct {
	Name                string `json:"name"`
	BaseAsset           string `json:"base_asset"`
	MakerFeeRateBase    string `json:"maker_fee_rate_base"`
	TakerFeeRateBase    string `json:"taker_fee_rate_base"`
	MakerFeeRateQuote   string `json:"maker_fee_rate_quote"`
	TakerFeeRateQuote   string `json:"taker_fee_rate_quote"`
	UnitAmount          string `json:"unit_amount"`
	LimitMaxAmount      string `json:"limit_max_amount"`
	MarketMaxAmount     string `json:"market_max_amount"`
	MarketAllowanceRate string `json:"market_allowance_rate"`
	PriceDigits         int    `json:"price_digits"`
	AmountDigits        int    `json:"amount_digits"`
	IsEnabled           bool   `json:"is_enabled"`
	StopOrder           bool   `json:"stop_order"`
	StopOrderAndCancel  bool   `json:"stop_order_and_cancel"`
}

// GetPairs is ...
func (api *APIClient) GetPairs() ([]Pairs, error) {
	url := "spot/pairs"
	resp, err := api.doRequest(restEndpoint, "GET", url, map[string]string{}, nil)
	log.Printf("action=Pairs.GetPairs, url=%s, resp=%s", url, string(resp))
	if err != nil {
		log.Printf("action=Pairs.GetPairs, err=%s", err.Error())
		return nil, err
	}
	var pairs []Pairs
	err = json.Unmarshal(resp, &pairs)
	if err != nil {
		log.Printf("action=Pairs.GetPairs, err=%s", err.Error())
		return nil, err
	}
	return pairs, nil
}

// TickerResponseInfo is ...
type TickerResponseInfo struct {
	Success int    `json:"success"`
	Data    Ticker `json:"data"`
}

// Ticker is ...
type Ticker struct {
	Sell      string `json:"sell"`
	Buy       string `json:"buy"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Last      string `json:"last"`
	Vol       string `json:"vol"`
	Timestamp int64  `json:"timestamp"`
}

// GetMiddle is ...
func (t *Ticker) GetMiddle() float64 {
	high, err := strconv.ParseFloat(t.High, 64)
	if err != nil {
		fmt.Printf("action=Ticker.GetMiddle, ticker=%v, err=%s\n", t, err.Error())
	}
	low, err := strconv.ParseFloat(t.Low, 64)
	if err != nil {
		fmt.Printf("action=Ticker.GetMiddle, ticker=%v, err=%s\n", t, err.Error())
	}
	middle := (high + low) / 2
	log.Printf("action=Ticker.GetMiddle, middle=%v", middle)
	return middle
}

// DateTime is ...
// func (t *Ticker) DateTime() time.Time {
// 	dateTimeUnix := time.Unix(t.Timestamp, 0)
// 	dateTimeStr := dateTimeUnix.Format(time.RFC3339)
// 	dateTime, err := time.Parse(time.RFC3339, dateTimeStr)
// 	if err != nil {
// 		log.Printf("action=DateTime, err=%s", err.Error())
// 	}
// 	return dateTime
// }

// TruncateDateTime is ...
// func (t *Ticker) TruncateDateTime(duration time.Duration) time.Time {
// 	return t.DateTime().Truncate(duration)
// }

// GetTicker is ...
func (api *APIClient) GetTicker(pair string) (*Ticker, error) {
	url := pair + "/ticker"
	resp, err := api.doRequest(publicEndpoint, "GET", url, map[string]string{"pair": pair}, nil)
	log.Printf("action=Ticker.GetTicker, url=%s, resp=%s", url, string(resp))
	if err != nil {
		log.Printf("action=Ticker.GetTicker, err=%s", err.Error())
		return nil, err
	}
	var tickerResponseInfo TickerResponseInfo
	err = json.Unmarshal(resp, &tickerResponseInfo)
	if err != nil {
		log.Printf("action=Ticker.GetTicker, err=%s", err.Error())
		return nil, err
	}
	return &tickerResponseInfo.Data, nil
}
