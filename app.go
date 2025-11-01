package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/ilius/price-widget/pkg/asset"
	"github.com/ilius/price-widget/pkg/config"
	"github.com/ilius/price-widget/pkg/managers/cryptomanager"
	"github.com/ilius/price-widget/pkg/managers/metals"
	qt "github.com/mappu/miqt/qt6"
)

const letfOrMiddleButton = qt.LeftButton | qt.MiddleButton

func Run() {
	qt.NewQApplication(os.Args)

	conf := config.Load()
	slog.Info("Using config", "conf", conf)

	allAssets := conf.Assets

	cryptoAssets := []*asset.Asset{}
	metalAssets := []*asset.Asset{}
	for _, asset := range allAssets {
		switch asset.Type {
		case "metal", "gold", "goldprice":
			metalAssets = append(metalAssets, asset)
		case "", "coin", "crypto", "coingecko":
			cryptoAssets = append(cryptoAssets, asset)
		}
	}
	refreshInterval := time.Duration(conf.RefreshIntervalSeconds) * time.Second

	cryptoManager := cryptomanager.New(cryptoAssets, refreshInterval)
	metalManager := metals.New(metalAssets, refreshInterval)

	getPrice := func(asset *asset.Asset) (float64, bool) {
		price, ok := metalManager.GetPrice(asset)
		if ok {
			return price, true
		}
		price, ok = cryptoManager.GetPrice(asset)
		if ok {
			return price, true
		}
		slog.Error("asset not found", "id", asset.ID, "name", asset.Name)
		return 0, false
	}

	// Create frameless black window
	window := qt.NewQWidget2()
	flag := qt.FramelessWindowHint | qt.Tool | qt.WindowStaysOnBottomHint
	if conf.BypassWindowManager {
		flag |= qt.BypassWindowManagerHint
	}
	window.SetWindowFlags(flag)

	window.SetStyleSheet("background-color: black;")

	// Layout setup
	rootLayout := qt.NewQHBoxLayout2()
	rootLayout.SetSpacing(30)
	rootLayout.SetContentsMargins(20, 20, 20, 20)

	cryptoManager.Init()
	metalManager.Init()

	font := qt.NewQFont()
	if conf.TextSize > 0 {
		font.SetPixelSize(conf.TextSize)
	}

	priceLabels := map[string]*qt.QLabel{}
	for _, asset := range allAssets {
		colLayout := qt.NewQVBoxLayout2()
		colLayout.SetSpacing(8)
		price, _ := getPrice(asset)
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

	showPrice := func(asset *asset.Asset, price float64) {
		label := priceLabels[asset.ID]
		label.SetText(asset.FormatPrice(price))
	}

	go cryptoManager.FetchLoop(showPrice)
	go metalManager.FetchLoop(showPrice)

	window.SetLayout(rootLayout.QLayout)
	window.Resize(len(allAssets)*100, 1)

	// --- Make it draggable ---
	var dragRelativePos *qt.QPoint

	exit := func() {
		// app.onExit()
		os.Exit(0)
	}

	actions := []*qt.QAction{}
	{
		action := qt.NewQAction2("Quit")
		action.OnTriggered(exit)
		actions = append(actions, action)
	}
	{
		action := qt.NewQAction2("Config")
		action.OnTriggered(func() {
			config.OpenInEditor(conf)
		})
		actions = append(actions, action)
	}
	{
		action := qt.NewQAction2("About")
		action.OnTriggered(func() {
			widget := qt.NewQDialog(window)
			aboutClickedWidget(widget.QWidget, nil)
			widget.Show()
		})
		actions = append(actions, action)
	}
	popupMenu := func(event *qt.QMouseEvent) {
		menu := qt.NewQMenu2()
		for _, action := range actions {
			menu.AddAction(action)
		}
		// menu.SetFont(app.systemDefaultFont)
		menu.Popup(event.GlobalPos())
	}

	window.OnMousePressEvent(func(super func(event *qt.QMouseEvent), event *qt.QMouseEvent) {
		// Note: event.Buttons() == event.Button()
		// slog.Info("OnMousePressEvent", "button", event.Button(), "buttons", event.Buttons())
		if event.Buttons()&letfOrMiddleButton > 0 {
			if os.Getenv("WAYLAND_DISPLAY") != "" {
				window.WindowHandle().StartSystemMove()
			} else {
				dragRelativePos = event.WindowPos().ToPoint()
			}
		} else if event.Buttons()&qt.RightButton > 0 {
			popupMenu(event)
		}
	})
	window.OnMouseMoveEvent(func(super func(event *qt.QMouseEvent), event *qt.QMouseEvent) {
		// Note: event.Button() is 0, but event.Buttons() is set (1 for left only).
		// slog.Info("OnMouseMoveEvent", "button", event.Button(), "buttons", event.Buttons())
		if dragRelativePos != nil && event.Buttons()&letfOrMiddleButton > 0 {
			window.Move(
				event.GlobalX()-dragRelativePos.X(),
				event.GlobalY()-dragRelativePos.Y(),
			)
		}
	})
	window.OnMouseReleaseEvent(func(super func(event *qt.QMouseEvent), event *qt.QMouseEvent) {
		dragRelativePos = nil
	})
	window.Show()
	qt.QApplication_Exec()
}
