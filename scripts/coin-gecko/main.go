// Takes a list of supported tokens on Polygon and formats them into a JSON file.
// Usage: go run scripts/coin-gecko/main.go

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const endpoint = "https://tokens.coingecko.com/polygon-pos/all.json"

type Token struct {
	ChainID  int    `json:"chainId"`
	Address  string `json:"address"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	LogoURI  string `json:"logoURI"`
}

func main() {
	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to fetch data. Status code:", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	tokens, ok := data["tokens"].([]interface{})
	if !ok {
		fmt.Println("No tokens found in the response.")
		return
	}

	formattedTokens := make(map[string]map[string]interface{})

	for _, token := range tokens {
		tok, ok := token.(map[string]interface{})
		if !ok {
			continue
		}

		symbol, ok := tok["symbol"].(string)
		if !ok {
			continue
		}

		address, ok := tok["address"].(string)
		if !ok {
			continue
		}

		decimals, ok := tok["decimals"].(float64)
		if !ok {
			continue
		}

		formattedTokens[symbol] = map[string]interface{}{
			"address":  address,
			"decimals": decimals,
		}
	}

	outputFile := "supportedTokens.json"
	outputJSON, err := json.MarshalIndent(formattedTokens, "", "    ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	err = os.WriteFile(outputFile, outputJSON, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Printf("Formatted tokens saved to %s\n", outputFile)
}
