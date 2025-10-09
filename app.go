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
	}

	// Create frameless black window
	window := qt.NewQWidget2()
	window.SetWindowFlags(qt.FramelessWindowHint | qt.Window)
	window.SetStyleSheet("background-color: black;")

	// Layout setup
	rootLayout := qt.NewQHBoxLayout2()
	rootLayout.SetSpacing(30)
	rootLayout.SetContentsMargins(20, 20, 20, 20)

	now := time.Now()
	var prices map[string]float64

	refreshInterval := time.Duration(conf.RefreshIntervalSeconds) * time.Second

	if cache != nil && now.Sub(cache.Timestamp) < refreshInterval {
		prices = cache.Prices
	} else {
		p, err := fetchCryptoPrices(assets)
		if err == nil {
			prices = p
			saveCache(&PriceCache{Timestamp: now, Prices: prices})
		} else if cache != nil {
			slog.Error("failed to fetch, using cached data", "err", err)
			prices = cache.Prices
		}
	}

	go func() {
		for {
			time.Sleep(refreshInterval)
			p, err := fetchCryptoPrices(assets)
			if err != nil {
				slog.Error("error fetching", "err", err)
				continue
			}
			saveCache(&PriceCache{Timestamp: time.Now(), Prices: p})
			// for _, asset := range assets {
			// 	if v, ok := p[asset.ID]; ok {
			// 		label := labels[asset.ID]
			// 		label.Text = fmt.Sprintf("%.*f", asset.Digits, v)
			// 		label.Refresh()
			// 	}
			// }
		}
	}()

	for _, asset := range assets {
		colLayout := qt.NewQVBoxLayout2()
		colLayout.SetSpacing(8)
		price := prices[asset.ID]
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
		}
		rootLayout.AddLayout(colLayout.QLayout)
		rootLayout.AddStretch()
	}

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
