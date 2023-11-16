package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const pubKey = "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEhqhj8rWPzkghzOZTUCOo/sdkE53sU1coVhaYskKGKrgiUF7lsSmxy46i3j8w7E7KMTfYBpCGAFYiWWARa0KQwg=="
const depthSize = 5

var HOST = "localhost"

type Ticker struct {
	Price  string `json:"price"`
	Symbol string `json:"symbol"`
}
type AddOrderReq struct {
	Price         string `json:"price"`
	Size          string `json:"size"`
	Side          string `json:"side"`
	Symbol        string `json:"symbol"`
	ClientOrderId string `json:"clientOrderId"`
}

func onTick(url string) *Ticker {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("HTTP request failed with status code %d\n", response.StatusCode)
		return nil
	}

	// Create a decoder for the response body
	var ticker Ticker
	decoder := json.NewDecoder(response.Body)

	// Decode the JSON response
	if err := decoder.Decode(&ticker); err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return nil
	}

	// Extract and print the price
	fmt.Printf("ETH-USD Price: %s\n", ticker.Price)
	return &ticker
}

func cancelAllOrders() {
	url := fmt.Sprintf("%s/api/v1/orders", HOST)

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new DELETE request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Add("X-Public-Key", pubKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Canceled all orders")
}

func placeOrder(side string, price, size decimal.Decimal) {
	client := &http.Client{}

	cOId := uuid.NewString()

	body := AddOrderReq{
		Price:         price.String(),
		Size:          size.String(),
		Side:          side,
		Symbol:        "ETH-USD",
		ClientOrderId: cOId,
	}
	url := fmt.Sprintf("%s/api/v1/order", HOST)
	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("error marshaling: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("error creating request: %v", err)
	}

	req.Header.Add("X-Public-Key", pubKey)

	res, err := client.Do(req)
	fmt.Printf("res is ------->: %#v\n", res)
	if err != nil {
		log.Fatalf("error post: %v", err)
	}

	fmt.Println("Created order with clientOrderId: ", cOId)
	fmt.Println("Status code:", res.StatusCode)

	defer res.Body.Close()
}

func updateOrders(price decimal.Decimal) {
	cancelAllOrders()
	factor := decimal.NewFromFloat(1.01)
	curPrice := price
	for i := 0; i < depthSize; i++ {
		curPrice = curPrice.Mul(factor)
		curSize := decimal.NewFromFloat(float64(i+1) * 10)
		placeOrder("sell", curPrice, curSize)
	}

	factor = decimal.NewFromFloat(0.99)
	curPrice = price
	for i := 0; i < depthSize; i++ {
		curPrice := curPrice.Mul(factor)
		curSize := decimal.NewFromFloat(float64(i+1) * 10)
		placeOrder("buy", curPrice, curSize)
	}

}
func main() {
	url := "https://www.binance.com/api/v3/ticker/price?symbol=ETHUSDT"
	println("Ticker URL: ", url)
	host := os.Getenv("ORDERBOOK_HOST")
	if len(host) > 0 {
		HOST = host
	}

	for {
		// Fetch the ticker price for ETH-USD
		ticker := onTick(url)
		if ticker != nil {
			price := decimal.RequireFromString(ticker.Price)
			updateOrders(price)

			fmt.Println("Sleeping for 10 seconds...")
			time.Sleep(10 * time.Second)
		}
	}
}
