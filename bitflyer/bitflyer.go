package bitflyer

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

	"github.com/gorilla/websocket"
)

const baseURL = "https://api.bitflyer.com/v1/"

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
	message := timestamp + method + endpoint + string(body)

	mac := hmac.New(sha256.New, []byte(api.secret))
	mac.Write([]byte(message))
	sign := hex.EncodeToString(mac.Sum(nil))

	return map[string]string{
		"ACCESS-KEY":       api.key,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      sign,
		"Content-Type":     "application/json",
	}
}

func (api *APIClient) doRequest(method, urlPath string, query map[string]string, data []byte) (body []byte, err error) {
	baseURL, err := url.Parse(baseURL)
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

// Balance is ...
type Balance struct {
	CurrentCode string  `json:"currency_code"`
	Amount      float64 `json:"amount"`
	Available   float64 `json:"available"`
}

// GetBalance is ...
func (api *APIClient) GetBalance() ([]Balance, error) {
	url := "me/getbalance"
	resp, err := api.doRequest("GET", url, map[string]string{}, nil)
	log.Printf("action=Balance.GetBalance, url=%s, resp=%s", url, string(resp))
	if err != nil {
		log.Printf("action=Balance.GetBalance, err=%s", err.Error())
		return nil, err
	}
	var balance []Balance
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		log.Printf("action=Balance.GetBalance, err=%s", err.Error())
		return nil, err
	}
	return balance, nil
}

// Ticker is ...
type Ticker struct {
	ProductCode     string  `json:"product_code"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

// GetMidPrice is ...
func (t *Ticker) GetMidPrice() float64 {
	midPrice := (t.BestBid + t.BestAsk) / 2
	log.Printf("action=Ticker.GetMidPrice, midPrice=%v", midPrice)
	return midPrice
}

// DateTime is ...
func (t *Ticker) DateTime() time.Time {
	dateTime, err := time.Parse(time.RFC3339, t.Timestamp)
	if err != nil {
		log.Printf("action=Ticker.DateTime, err=%s", err.Error())
	}
	log.Printf("action=Ticker.DateTime, dateTime=%v", dateTime)
	return dateTime
}

// TruncateDateTime is ...
func (t *Ticker) TruncateDateTime(duration time.Duration) time.Time {
	truncateDateTime := t.DateTime().Truncate(duration)
	log.Printf("action=Ticker.TruncateDateTime, truncateDateTime=%v", truncateDateTime)
	return truncateDateTime
}

// GetTicker is ...
func (api *APIClient) GetTicker(productCode string) (*Ticker, error) {
	url := "ticker"
	resp, err := api.doRequest("GET", url, map[string]string{"product_code": productCode}, nil)
	log.Printf("action=APIClient.GetTicker, url=%s, resp=%s", url, string(resp))
	if err != nil {
		log.Printf("action=APIClient.GetTicker, err=%s", err.Error())
		return nil, err
	}
	var ticker Ticker
	err = json.Unmarshal(resp, &ticker)
	if err != nil {
		log.Printf("action=APIClient.GetTicker, err=%s", err.Error())
		return nil, err
	}
	return &ticker, nil
}

// JSONRPC2 is ...
type JSONRPC2 struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result,omitempty"`
	ID      *int        `json:"id,omitempty"`
}

// SubscribeParams is ...
type SubscribeParams struct {
	Channel string `json:"channel"`
}

// GetRealTimeTicker is ...
func (api *APIClient) GetRealTimeTicker(pair string, ch chan<- Ticker) {
	u := url.URL{Scheme: "wss", Host: "ws.lightstream.bitflyer.com", Path: "/json-rpc"}
	log.Printf("action=APIClient.GetRealTimeTicker, connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	channel := fmt.Sprintf("lightning_ticker_%s", pair)
	if err := c.WriteJSON(&JSONRPC2{Version: "2.0", Method: "subscribe", Params: &SubscribeParams{channel}}); err != nil {
		log.Fatal("subscribe:", err)
		return
	}

OUTER:
	for {
		message := new(JSONRPC2)
		if err := c.ReadJSON(message); err != nil {
			log.Println("read:", err)
			return
		}

		if message.Method == "channelMessage" {
			switch v := message.Params.(type) {
			case map[string]interface{}:
				for key, binary := range v {
					if key == "message" {
						marshaTic, err := json.Marshal(binary)
						if err != nil {
							continue OUTER
						}
						var ticker Ticker
						if err := json.Unmarshal(marshaTic, &ticker); err != nil {
							continue OUTER
						}
						ch <- ticker
					}
				}
			}
		}
	}
}

// Order is ...
type Order struct {
	ID                     int     `json:"id"`
	ChildOrderAcceptanceID string  `json:"child_order_acceptance_id"`
	ProductCode            string  `json:"product_code"`
	ChildOrderType         string  `json:"child_order_type"`
	Side                   string  `json:"side"`
	Price                  float64 `json:"price"`
	Size                   float64 `json:"size"`
	MinuteToExpires        int     `json:"minute_to_expire"`
	TimeInForce            string  `json:"time_in_force"`
	Status                 string  `json:"status"`
	ErrorMessage           string  `json:"error_message"`
	AveragePrice           float64 `json:"average_price"`
	ChildOrderState        string  `json:"child_order_state"`
	ExpireDate             string  `json:"expire_date"`
	ChildOrderDate         string  `json:"child_order_date"`
	OutstandingSize        float64 `json:"outstanding_size"`
	CancelSize             float64 `json:"cancel_size"`
	ExecutedSize           float64 `json:"executed_size"`
	TotalCommission        float64 `json:"total_commission"`
	Count                  int     `json:"count"`
	Before                 int     `json:"before"`
	After                  int     `json:"after"`
}

// ResponseSendChildOrder is ...
type ResponseSendChildOrder struct {
	ChildOrderAcceptanceID string `json:"child_order_acceptance_id"`
}

// SendOrder is ...
func (api *APIClient) SendOrder(order *Order) (*ResponseSendChildOrder, error) {
	data, err := json.Marshal(order)
	if err != nil {
		log.Printf("action=APIClient.SendOrder, err=%s", err.Error())
		return nil, err
	}
	url := "me/sendchildorder"
	resp, err := api.doRequest("POST", url, map[string]string{}, data)
	log.Printf("action=APIClient.SendOrder, url=%s, resp=%s", url, string(resp))
	if err != nil {
		log.Printf("action=APIClient.SendOrder, err=%s", err.Error())
		return nil, err
	}
	var response ResponseSendChildOrder
	err = json.Unmarshal(resp, &response)
	if err != nil {
		log.Printf("action=APIClient.SendOrder, err=%s", err.Error())
		return nil, err
	}
	return &response, nil
}

// GetOrders is ...
func (api *APIClient) GetOrders(query map[string]string) ([]Order, error) {
	url := "me/getchildorders"
	resp, err := api.doRequest("GET", url, query, nil)
	log.Printf("action=Ticker.GetOrders, url=%s, resp=%s", url, string(resp))
	if err != nil {
		log.Printf("action=APIClient.GetOrders, err=%s", err.Error())
		return nil, err
	}
	var response []Order
	err = json.Unmarshal(resp, &response)
	if err != nil {
		log.Printf("action=APIClient.GetOrders, err=%s", err.Error())
		return nil, err
	}
	return response, nil
}
