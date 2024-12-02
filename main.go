package main

import (
	"log"
	"os"
	"sicsimgo/internal"
	"sicsimgo/ui"

	"gioui.org/app"
	"gioui.org/unit"
)

func main() {
	go func() {

		// Create a window
		w := new(app.Window)
		internal.ResetWIndowTitle(w)
		w.Option(app.Size(unit.Dp(1050), unit.Dp(450)))

		if err := ui.DrawWindow(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
