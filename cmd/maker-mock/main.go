package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Ticker struct {
	Price  string `json:"price"`
	Symbol string `json:"symbol"`
}

func doit(url string) *Ticker {
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

func main() {

	url := "https://www.binance.com/api/v3/ticker/price?symbol=ETHUSDT"
	for {
		// Fetch the ticker price for ETH-USD
		ticker := doit(url)
		if ticker != nil {
			fmt.Printf("ETH-USD Price: %s\n", ticker.Price)

			// Sleep for 10 seconds
			time.Sleep(10 * time.Second)
		}
	}
}
