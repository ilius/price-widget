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
	Currency    string  `json:"curr"`     // "USD"
	GoldPrice   float64 `json:"xauPrice"` // Gold price per ounce
	SilverPrice float64 `json:"xagPrice"` // Silver price per ounce

	// GoldPriceClose   float64 `json:"xauClose"` // Gold price per ounce at the close of the previous trading session
	// SilverPriceClose float64 `json:"xagClose"` // Silver price per ounce at the close of the previous trading session
}

// Example of all fields in goldApiResponseItem:
// 		"chgXag": -0.1249,
// 		"chgXau": 0.325,
// 		"curr": "USD",
// 		"pcXag": -0.2559,
// 		"pcXau": 0.0081,
// 		"xagClose": 48.80345,
// 		"xagPrice": 48.6785,
// 		"xauClose": 4002.605,
// 		"xauPrice": 4002.93
// chg (prefix): Absolute change in price
// pc (prefix): Percent change in price

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

	goldUsdOz := data.Items[0].GoldPrice       // USD per oz
	goldUsdKg := goldUsdOz / troyOunceToKg     // USD per kg
	silverUsdOz := data.Items[0].SilverPrice   // USD per oz
	silverUsdKg := silverUsdOz / troyOunceToKg // USD per kg

	return map[string]float64{
		"gold":      goldUsdOz,
		"gold_oz":   goldUsdOz,
		"gold_kg":   goldUsdKg,
		"silver":    silverUsdOz,
		"silver_oz": silverUsdOz,
		"silver_kg": silverUsdKg,
	}, nil
}
