package coingecko

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/ilius/price-widget/pkg/asset"
)

func New() *provider {
	return &provider{}
}

type provider struct{}

func (*provider) SupportedIDs() map[string]struct{} {
	return nil
}

func (*provider) FetchPrices(assets []*asset.Asset) (map[string]float64, error) {
	slog.Info("fetching prices...")

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", idsToCSV(assets))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]map[string]float64
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	prices := make(map[string]float64)
	for _, asset := range assets {
		if v, ok := data[asset.ID]["usd"]; ok {
			prices[asset.ID] = v
		}
	}
	slog.Info("fetched prices")

	return prices, nil
}

func idsToCSV(assets []*asset.Asset) string {
	out := ""
	for i, asset := range assets {
		if i > 0 {
			out += ","
		}
		out += asset.ID
	}
	return out
}
