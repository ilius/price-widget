package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const troyOunceToKg = 0.0311035

// This endpoint is semi-public and designed to power the
// goldprice.org live charts and widgets.
// It’s not an official API, but people (and even trading dashboards)
// have used it reliably for years.
// Auth: None required (open endpoint)
// Rate limit: Not formally documented — generally tolerant
// Safe polling interval: 30 to 60 seconds is considered safe
// Risk of blocking: Very low if you keep requests below ~60/hour
// Response size: Small (a few hundred bytes)
// Update frequency: Roughly every 60 seconds
const goldApiUrl = "https://data-asg.goldprice.org/dbXRates/USD"

type goldApiResponseItem struct {
	Currency  string  `json:"curr"`     // "USD"
	XAU_Ounce float64 `json:"xauPrice"` // Gold price per ounce
	XAG_Ounce float64 `json:"xagPrice"` // Silver price per ounce
}

type goldApiResponse struct {
	Items []goldApiResponseItem `json:"items"`
}

func fetchPrices() (*goldApiResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(goldApiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data goldApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func main() {
	data, err := fetchPrices()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(data.Items) == 0 {
		fmt.Println("No data returned")
		return
	}

	// data.Items[0].Currency == "USD"

	gold := data.Items[0].XAU_Ounce
	silver := data.Items[0].XAG_Ounce

	goldPerKg := gold / troyOunceToKg
	silverPerKg := silver / troyOunceToKg

	fmt.Printf("Gold:  \t\t%.2f USD/oz\t\t%.2f USD/kg\n", gold, goldPerKg)
	fmt.Printf("Silver:\t\t%.2f USD/oz\t\t%.2f USD/kg\n", silver, silverPerKg)
}
