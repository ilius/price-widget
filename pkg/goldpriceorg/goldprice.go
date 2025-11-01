package goldpriceorg

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ilius/price-widget/pkg/asset"
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

func New() *provider {
	return &provider{}
}

type provider struct{}

func (*provider) SupportedIDs() map[string]struct{} {
	return map[string]struct{}{
		"gold":      {},
		"gold_oz":   {},
		"gold_kg":   {},
		"silver":    {},
		"silver_oz": {},
		"silver_kg": {},
	}
}

func (*provider) FetchPrices(assets []*asset.Asset) (map[string]float64, error) {
	// assetIds := asset.IdMap(assets)
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
	// data.Items[0].Currency == "USD"

	goldUsdOz := data.Items[0].XAU_Ounce       // USD per oz
	silverUsdOz := data.Items[0].XAG_Ounce     // USD per oz
	goldUsdKg := goldUsdOz / troyOunceToKg     // USD per kg
	silverUsdKg := silverUsdOz / troyOunceToKg // USD per kg

	prices := map[string]float64{}
	for _, asset := range assets {
		switch asset.ID {
		case "gold", "gold_oz":
			prices[asset.ID] = goldUsdOz
		case "gold_kg":
			prices[asset.ID] = goldUsdKg
		case "silver", "silver_oz":
			prices[asset.ID] = silverUsdOz
		case "silver_kg":
			prices[asset.ID] = silverUsdKg
		}
	}

	return prices, nil
}
