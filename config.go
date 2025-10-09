package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Config file paths
var (
	configDir      = filepath.Join(os.Getenv("HOME"), ".config", "price-widget")
	priceCacheFile = filepath.Join(configDir, "prices.json")
	configFile     = filepath.Join(configDir, "config.toml")
)

// TOML config structure
type Config struct {
	TextSize float32 `toml:"text_size"`

	RefreshIntervalSeconds int `toml:"refresh_interval_seconds"`

	Assets []*Asset `toml:"assets"`
}

func loadConfig() *Config {
	conf := &Config{
		TextSize:               24,
		RefreshIntervalSeconds: 15 * 60,
		Assets:                 []*Asset{},
	}

	b, err := os.ReadFile(configFile)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Error("error in reading config toml file", "err", err)
		} else {
			slog.Warn("config toml file not found", "configFile", configFile)
		}
		return conf
	}

	err = toml.Unmarshal(b, conf)
	if err != nil {
		slog.Error("error in loading config toml file", "err", err)
	}
	return conf
}
