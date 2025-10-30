package main

import (
	"fmt"
	"log/slog"
	"runtime"

	qt "github.com/mappu/miqt/qt6"
)

var isMac = runtime.GOOS == "darwin"

func addTabWithIcon(
	tabWidget *qt.QTabWidget,
	widget *qt.QWidget,
	label string,
	filename string,
) {
	if isMac {
		_ = tabWidget.AddTab(widget, label)
		return
	}
	icon, err := loadPNGIcon(filename)
	if err != nil {
		slog.Error("error loading icon", "filename", filename)
	}
	if icon == nil {
		_ = tabWidget.AddTab(widget, label)
		return
	}
	_ = tabWidget.AddTab2(widget, icon, label)
}

func aboutClickedWidget(widget *qt.QWidget, icon *qt.QIcon) {
	widget.SetWindowTitle("About " + APP_DESC)
	widget.Resize(700, 500)
	if icon != nil {
		widget.SetWindowIcon(icon)
	}

	topHBox := qt.NewQFrame(nil)
	topHBoxLayout := qt.NewQHBoxLayout(topHBox.QWidget)

	// {
	// 	pixmap, err := loadPNGPixmap(".png")
	// 	if err != nil {
	// 		slog.Error("failed to load icon image", "err", err)
	// 	} else {
	// 		label := qt.NewQLabel2()
	// 		label.SetPixmap(pixmap)
	// 		label.SetMinimumWidth(80)
	// 		topHBoxLayout.AddWidget3(label.QWidget, 0, qt.AlignCenter)
	// 	}
	// }

	topLabel := qt.NewQLabel3(fmt.Sprintf(
		"%s version %s\nUsing Qt %v and Go %v",
		APP_DESC,
		VERSION,
		qt.QLibraryInfo_Version().ToString(),
		runtime.Version()[2:],
	))
	topHBoxLayout.AddWidget(topLabel.QWidget)
	topHBoxLayout.AddStretch()

	tabWidget := qt.NewQTabWidget2()
	tabWidget.SetSizePolicy2(expanding, expanding)
	tabWidget.SetIconSize(qt.NewQSize2(22, 22))

	// tabWidget.SetTabPosition(qt.QTabWidget__West)
	// tabBar := tabWidget.TabBar()
	// tabWidget.SetStyleSheet(`
	// QTabBar::tab {
	// 	direction: ltr;
	// 	padding: 15px;
	// }`)

	aboutLabel := qt.NewQLabel3(ABOUT)
	aboutLabel.SetTextInteractionFlags(qt.TextSelectableByMouse)
	aboutLabel.SetAlignment(qt.AlignTop)
	aboutLabel.SetWordWrap(true)
	addTabWithIcon(tabWidget, aboutLabel.QWidget, "About", "dialog-information-22.png")

	authorsLabel := qt.NewQLabel3(AUTHORS)
	authorsLabel.SetTextInteractionFlags(qt.TextSelectableByMouse)
	authorsLabel.SetAlignment(qt.AlignTop)
	addTabWithIcon(tabWidget, authorsLabel.QWidget, "Authors", "author-22.png")

	licenseWidget := qt.NewQTextEdit2()
	licenseWidget.SetReadOnly(true)
	licenseWidget.SetPlainText(LICENSE)
	addTabWithIcon(tabWidget, licenseWidget.QWidget, "License", "license-22.png")

	buttonBox := qt.NewQDialogButtonBox2()
	buttonBox.AddButton2("Website", qt.QDialogButtonBox__AcceptRole).OnClicked(func() {
		qt.QDesktopServices_OpenUrl(qt.NewQUrl3(WEBSITE))
	})
	closeButton := buttonBox.AddButton2("  Close  ", qt.QDialogButtonBox__AcceptRole)
	closeButton.OnClicked(func() {
		_ = widget.Close()
	})
	// closeButton.SetDefault(true)

	mainBox := qt.NewQVBoxLayout(widget)
	mainBox.AddWidget(topHBox.QWidget)
	mainBox.AddWidget(tabWidget.QWidget)
	mainBox.AddWidget(buttonBox.QWidget)
}
