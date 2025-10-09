package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	qt "github.com/mappu/miqt/qt6"
)

func Run() {
	qt.NewQApplication(os.Args)

	conf := loadConfig()
	slog.Info("Using config", "spacing", conf.Spacing, "text_size", conf.TextSize)

	assets := conf.Assets

	cache, err := loadPriceCache()
	if err != nil {
		slog.Error("error loading price cache", "err", err)
		cache = &PriceCache{Prices: map[string]float64{}}
	}

	// Create frameless black window
	window := qt.NewQWidget2()
	window.SetWindowFlags(qt.FramelessWindowHint | qt.Window)
	window.SetStyleSheet("background-color: black;")

	// Layout setup
	rootLayout := qt.NewQHBoxLayout2()
	rootLayout.SetSpacing(30)
	rootLayout.SetContentsMargins(20, 20, 20, 20)

	refreshInterval := time.Duration(conf.RefreshIntervalSeconds) * time.Second

	now := time.Now()
	if now.Sub(cache.Timestamp) > refreshInterval {
		prices, err := fetchCryptoPrices(assets)
		if err != nil {
			slog.Error("failed to fetch, using cached data", "err", err)

		}
		cache.Prices = prices
		cache.Timestamp = now
		saveCache(cache)
	}

	priceLabels := map[string]*qt.QLabel{}
	for _, asset := range assets {
		colLayout := qt.NewQVBoxLayout2()
		colLayout.SetSpacing(8)
		price := cache.Prices[asset.ID]
		priceStr := fmt.Sprintf("%.*f", asset.Digits, price)
		{
			label := qt.NewQLabel3(asset.Name)
			label.SetAlignment(qt.AlignCenter)
			label.SetStyleSheet("color: white; font-size: 22px;")
			colLayout.AddWidget(label.QWidget)
		}
		rootLayout.AddSpacing(10)
		{
			label := qt.NewQLabel3(priceStr)
			label.SetAlignment(qt.AlignCenter)
			label.SetStyleSheet("color: white; font-size: 22px;")
			colLayout.AddWidget(label.QWidget)
			priceLabels[asset.ID] = label
		}
		rootLayout.AddLayout(colLayout.QLayout)
		rootLayout.AddStretch()
	}

	go func() {
		ticker := time.NewTicker(refreshInterval)
		for {
			<-ticker.C
			prices, err := fetchCryptoPrices(assets)
			if err != nil {
				slog.Error("error fetching", "err", err)
				continue
			}
			cache.Prices = prices
			cache.Timestamp = time.Now()
			saveCache(cache)
			for _, asset := range assets {
				if price, ok := prices[asset.ID]; ok {
					label := priceLabels[asset.ID]
					label.SetText(fmt.Sprintf("%.*f", asset.Digits, price))
				}
			}
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
