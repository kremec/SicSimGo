package internal

import "gioui.org/app"

var DefaultWindowTitle = "SicSimGo"

func SetWIndowTitle(title string, w *app.Window) {
	w.Option(app.Title(title))
}

func ResetWIndowTitle(w *app.Window) {
	w.Option(app.Title(DefaultWindowTitle))
}
