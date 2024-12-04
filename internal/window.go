package internal

import "gioui.org/app"

var DefaultWindowTitle = "SicSimGo"

func SetWindowTitle(title string, w *app.Window) {
	w.Option(app.Title(title))
}

func ResetWindowTitle(w *app.Window) {
	w.Option(app.Title(DefaultWindowTitle))
}
