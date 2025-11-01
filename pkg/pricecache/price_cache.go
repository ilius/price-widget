package pricecache

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func New() *PriceCache {
	return &PriceCache{Prices: map[string]float64{}}
}

type PriceCache struct {
	LastFetch time.Time          `json:"timestamp"`
	Prices    map[string]float64 `json:"prices"`
}

func (c *PriceCache) Load(fpath string) error {
	b, err := os.ReadFile(fpath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, c); err != nil {
		return err
	}
	return nil
}

func (c *PriceCache) Save(fpath string) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(fpath), 0o755)
	if err != nil {
		slog.Error("error creating directory", "err", err)
	}
	return os.WriteFile(fpath, b, 0o644)
}
