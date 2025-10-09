package main

import (
	"os"

	qt "github.com/mappu/miqt/qt6"
)

var table = [][]string{
	{"Bitcoin", "Litecoin", "Nano"},
	{"120000", "116", "0.821"},
}

func main() {
	qt.NewQApplication(os.Args)

	// Create frameless black window
	window := qt.NewQWidget2()
	window.SetWindowFlags(qt.FramelessWindowHint | qt.Window)
	window.SetStyleSheet("background-color: black;")

	// Layout setup
	rootLayout := qt.NewQHBoxLayout2()
	rootLayout.SetSpacing(30)
	rootLayout.SetContentsMargins(20, 20, 20, 20)

	numCols := len(table[0])
	for c := 0; c < numCols; c++ {
		colLayout := qt.NewQVBoxLayout2()
		colLayout.SetSpacing(8)

		for r := 0; r < len(table); r++ {
			label := qt.NewQLabel3(table[r][c])
			label.SetAlignment(qt.AlignCenter)
			label.SetStyleSheet("color: white; font-size: 22px;")
			colLayout.AddWidget(label.QWidget)
		}

		rootLayout.AddLayout(colLayout.QLayout)
	}

	window.SetLayout(rootLayout.QLayout)
	window.Resize(520, 160)

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
