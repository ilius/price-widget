package managers

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/ilius/price-widget/pkg/asset"
	"github.com/ilius/price-widget/pkg/config"
	"github.com/ilius/price-widget/pkg/pricecache"
	"github.com/ilius/price-widget/pkg/providers"
)

func NewManager(
	provider providers.Provider,
	assets []*asset.Asset,
	refresh time.Duration,
) *Manager {
	cacheFile := filepath.Join(config.Dir(), "prices-"+provider.Name()+".json")
	return &Manager{
		provider:  provider,
		assets:    assets,
		cache:     pricecache.New(),
		cacheFile: cacheFile,
		refresh:   refresh,
	}
}

type Manager struct {
	provider  providers.Provider
	assets    []*asset.Asset
	cache     *pricecache.PriceCache
	cacheFile string
	refresh   time.Duration
}

func (m *Manager) Init() {
	if len(m.assets) == 0 {
		return
	}
	err := m.cache.Load(m.cacheFile)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Error("error loading price cache", "err", err)
		}
	}
	now := time.Now()
	if now.Sub(m.cache.LastFetch) > m.refresh {
		slog.Info("manager: fetching", "cacheFile", m.cacheFile)
		prices, err := m.provider.FetchPrices(m.assets)
		if err != nil {
			slog.Error("failed to fetch, using cached data", "err", err)
		}
		m.cache.Prices = prices
		m.cache.LastFetch = now
		m.cache.Save(m.cacheFile)
	}
}

func (m *Manager) FetchAndSave() bool {
	if len(m.assets) == 0 {
		return true
	}
	prices, err := m.provider.FetchPrices(m.assets)
	if err != nil {
		slog.Error("error fetching", "err", err)
		return false
	}
	m.cache.Prices = prices
	m.cache.LastFetch = time.Now()
	m.cache.Save(m.cacheFile)
	return true
}

func (m *Manager) GetPrice(asset *asset.Asset) (float64, bool) {
	price, ok := m.cache.Prices[asset.ID]
	return price, ok
}

func (m *Manager) FetchLoop(showPrice func(*asset.Asset, float64)) {
	if len(m.assets) == 0 {
		return
	}
	now := time.Now()
	lastTime := m.cache.LastFetch.Truncate(time.Minute)
	sleepDuration := lastTime.Add(m.refresh).Sub(now)
	slog.Info("sleeping", "duration", sleepDuration, "last_time", lastTime, "now", now)
	time.Sleep(sleepDuration)
	ticker := time.NewTicker(m.refresh)
	for {
		m.FetchAndSave()
		for _, asset := range m.assets {
			price, ok := m.GetPrice(asset)
			if ok {
				showPrice(asset, price)
			}
		}
		<-ticker.C
	}
}
