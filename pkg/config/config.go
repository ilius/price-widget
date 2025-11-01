package config

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ilius/price-widget/pkg/asset"
	qt "github.com/mappu/miqt/qt6"
	"github.com/pelletier/go-toml/v2"
)

// Config file paths
var (
	configDir  = filepath.Join(os.Getenv("HOME"), ".config", "price-widget")
	configFile = filepath.Join(configDir, "config.toml")
)

func Dir() string {
	return configDir
}

// TOML config structure
type Config struct {
	TextSize int `toml:"text_size"`

	GridSpacing   int `toml:"grid_spacing"`
	WindowMargins int `toml:"window_margins"`

	RefreshIntervalSeconds int `toml:"refresh_interval_seconds"`

	BypassWindowManager bool `toml:"bypass_window_manager"`

	Assets []*asset.Asset `toml:"assets"`
}

func Load() *Config {
	conf := &Config{
		TextSize:               24,
		GridSpacing:            10,
		WindowMargins:          15,
		RefreshIntervalSeconds: 15 * 60,
		Assets: []*asset.Asset{
			{
				Name:           "Bitcoin",
				ID:             "bitcoin",
				Digits:         -1,
				HumanizeFormat: "###,###.",
			},
		},
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

func Save(conf *Config) error {
	err := os.MkdirAll(configDir, 0o755)
	if err != nil {
		return err
	}
	b, err := toml.Marshal(conf)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, b, 0o644)
}

func ensureConfigExists(conf *Config) bool {
	stat, err := os.Stat(configFile)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Error("error checking config file", "err", err)
			return false
		}
		err := Save(conf)
		if err != nil {
			slog.Error("error saving config file", "err", err)
			return false
		}
		return true
	}
	if stat.IsDir() {
		slog.Error("config file is a directory", "path", configFile)
		return false
	}
	return true
}

func OpenInEditor(conf *Config) {
	if !ensureConfigExists(conf) {
		return
	}
	url := qt.NewQUrl()
	url.SetScheme("file")
	url.SetPath2(configFile, qt.QUrl__TolerantMode)
	_ = qt.QDesktopServices_OpenUrl(url)
}
