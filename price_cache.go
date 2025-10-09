package main

import (
	"encoding/json"
	"os"
	"time"
)

type PriceCache struct {
	Timestamp time.Time          `json:"timestamp"`
	Prices    map[string]float64 `json:"prices"`
}

func loadPriceCache() (*PriceCache, error) {
	b, err := os.ReadFile(priceCacheFile)
	if err != nil {
		return nil, err
	}
	c := &PriceCache{}
	if err := json.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
}

func saveCache(data *PriceCache) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	os.MkdirAll(configDir, 0o755)
	return os.WriteFile(priceCacheFile, b, 0o644)
}
