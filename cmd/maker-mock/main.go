package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

const pubKey = "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEhqhj8rWPzkghzOZTUCOo"
const depthSize = 5

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
	url := "http://localhost/api/v1/orders"

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new DELETE request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add a custom header to the request
	req.Header.Add("X-Public-Key", pubKey) // Replace "YourAccessToken" with your actual token

	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
}

func placeOrder(side string, price, size decimal.Decimal) {
	req := AddOrderReq{
		Price:         price.String(),
		Size:          size.String(),
		Side:          side,
		Symbol:        "ETH_USD",
		ClientOrderId: "f677273e-12de-4acc-a4f8-de7fb5b86e37",
	}
	url := "http://localhost/api/v1/orders"
	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("error marshaling: %v", err)
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("error post: %v", err)
	}

	// Read the response body line by line
	defer res.Body.Close()

}
func updateOrders(price decimal.Decimal) {
	cancelAllOrders()
	factor := decimal.NewFromFloat(1.01)
	curPrice := price
	for i := 0; i < depthSize; i++ {
		curPrice = price.Mul(factor)
		curSize := decimal.NewFromFloat(float64(i+1) * 10)
		placeOrder("SELL", curPrice, curSize)

	}

	factor = decimal.NewFromFloat(0.99)
	curPrice = price
	for i := 0; i < depthSize; i++ {
		curPrice := curPrice.Mul(factor)
		curSize := decimal.NewFromFloat(float64(i+1) * 10)
		placeOrder("BUY", curPrice, curSize)

	}

	placeOrder("SELL", price, decimal.NewFromFloat(1.1))
}
func main() {

	url := "https://www.binance.com/api/v3/ticker/price?symbol=ETHUSDT"
	for {
		// Fetch the ticker price for ETH-USD
		ticker := onTick(url)
		if ticker != nil {
			fmt.Printf("ETH-USD Price: %s\n", ticker.Price)
			price := decimal.RequireFromString(ticker.Price)
			updateOrders(price)

			// Sleep for 10 seconds
			time.Sleep(10 * time.Second)
		}
	}
}
