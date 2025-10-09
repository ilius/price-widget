package main

import (
	"encoding/json"
	"os"
	"time"
)

type CryptoPriceCache struct {
	LastFetch time.Time          `json:"timestamp"`
	Prices    map[string]float64 `json:"prices"`
}

func loadCryptoPriceCache() (*CryptoPriceCache, error) {
	b, err := os.ReadFile(priceCacheFile)
	if err != nil {
		return nil, err
	}
	c := &CryptoPriceCache{}
	if err := json.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
}

func saveCryptoCache(data *CryptoPriceCache) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	os.MkdirAll(configDir, 0o755)
	return os.WriteFile(priceCacheFile, b, 0o644)
}
