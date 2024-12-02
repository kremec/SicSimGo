package events

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget/material"
)

func HandleGlobalEvents(gtx layout.Context, theme *material.Theme) {
	for {
		event, ok := gtx.Event(
			key.Filter{
				Name:     key.Name("+"),
				Optional: key.ModCtrl,
			},
			key.Filter{
				Name:     key.Name("-"),
				Optional: key.ModCtrl,
			},
		)
		if !ok {
			break
		}

		switch event := event.(type) {
		case key.Event:
			switch event.Name {
			case "+":
				if event.Modifiers.Contain(key.ModCtrl) && event.State == key.Press {
					theme.TextSize += 1
				}
			case "-":
				if event.Modifiers.Contain(key.ModCtrl) && event.State == key.Press {
					theme.TextSize -= 1
				}
			}
		}
	}
}
