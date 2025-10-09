package main

import (
	"log/slog"
	"os"
	"time"

	qt "github.com/mappu/miqt/qt6"
)

func Run() {
	qt.NewQApplication(os.Args)

	conf := loadConfig()
	slog.Info("Using config", "conf", conf)

	assets := conf.Assets

	cache, err := loadPriceCache()
	if err != nil {
		slog.Error("error loading price cache", "err", err)
		cache = &PriceCache{Prices: map[string]float64{}}
	}

	// Create frameless black window
	window := qt.NewQWidget2()
	window.SetWindowFlags(qt.FramelessWindowHint | qt.Tool | qt.WindowStaysOnBottomHint)

	window.SetStyleSheet("background-color: black;")

	// Layout setup
	rootLayout := qt.NewQHBoxLayout2()
	rootLayout.SetSpacing(30)
	rootLayout.SetContentsMargins(20, 20, 20, 20)

	refreshInterval := time.Duration(conf.RefreshIntervalSeconds) * time.Second

	now := time.Now()
	if now.Sub(cache.LastFetch) > refreshInterval {
		prices, err := fetchCryptoPrices(assets)
		if err != nil {
			slog.Error("failed to fetch, using cached data", "err", err)

		}
		cache.Prices = prices
		cache.LastFetch = now
		saveCache(cache)
	}

	font := qt.NewQFont()
	if conf.TextSize > 0 {
		font.SetPixelSize(conf.TextSize)
	}

	priceLabels := map[string]*qt.QLabel{}
	for _, asset := range assets {
		colLayout := qt.NewQVBoxLayout2()
		colLayout.SetSpacing(8)
		price := cache.Prices[asset.ID]
		{
			label := qt.NewQLabel3(asset.Name)
			label.SetAlignment(qt.AlignCenter)
			label.SetFont(font)
			label.SetStyleSheet("color: white;")
			colLayout.AddWidget(label.QWidget)
		}
		rootLayout.AddSpacing(10)
		{
			label := qt.NewQLabel3(asset.FormatPrice(price))
			label.SetAlignment(qt.AlignCenter)
			label.SetFont(font)
			label.SetStyleSheet("color: white;")
			colLayout.AddWidget(label.QWidget)
			priceLabels[asset.ID] = label
		}
		rootLayout.AddLayout(colLayout.QLayout)
		rootLayout.AddStretch()
	}

	fetch := func() bool {
		prices, err := fetchCryptoPrices(assets)
		if err != nil {
			slog.Error("error fetching", "err", err)
			return false
		}
		cache.Prices = prices
		cache.LastFetch = time.Now()
		saveCache(cache)
		for _, asset := range assets {
			if price, ok := prices[asset.ID]; ok {
				label := priceLabels[asset.ID]
				label.SetText(asset.FormatPrice(price))
			}
		}
		return true
	}

	go func() {
		now := time.Now()
		lastTime := cache.LastFetch.Truncate(time.Minute)
		sleepDuration := lastTime.Add(refreshInterval).Sub(now)
		slog.Info("sleeping", "duration", sleepDuration, "last_time", lastTime, "now", now)
		time.Sleep(sleepDuration)
		ticker := time.NewTicker(refreshInterval)
		for {
			fetch()
			<-ticker.C
		}
	}()

	window.SetLayout(rootLayout.QLayout)
	window.Resize(len(assets)*100, 1)

	// --- Make it draggable ---
	var dragRelativePos *qt.QPoint
	window.OnMousePressEvent(func(super func(event *qt.QMouseEvent), event *qt.QMouseEvent) {
		if event.Button() == qt.LeftButton {
			dragRelativePos = event.Pos()
		}
	})

	window.OnMouseMoveEvent(func(super func(event *qt.QMouseEvent), event *qt.QMouseEvent) {
		if event.Buttons()&qt.LeftButton != 0 {
			window.Move(
				event.GlobalX()-dragRelativePos.X(),
				event.GlobalY()-dragRelativePos.Y(),
			)
		}

	})
	window.Show()
	qt.QApplication_Exec()
}
