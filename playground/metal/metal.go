package main

import (
	"encoding/json"
	"fmt"

	"github.com/ilius/price-widget/pkg/asset"
	"github.com/ilius/price-widget/pkg/goldpriceorg"
)

func main() {
	provider := goldpriceorg.New()
	prices, err := provider.FetchPrices([]*asset.Asset{
		{ID: "gold_oz"},
		{ID: "gold_kg"},
		{ID: "silver_oz"},
		{ID: "silver_kg"},
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	pricesJson, err := json.MarshalIndent(prices, "", "    ")

	fmt.Println(string(pricesJson))
}
