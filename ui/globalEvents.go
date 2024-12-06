package ui

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget/material"
)

const debugHandleGlobalEvents bool = false

func HandleGlobalEvents(gtx layout.Context, theme *material.Theme, w *app.Window) {
	event, _ := gtx.Event(
		key.Filter{
			Name:     key.Name("+"),
			Optional: key.ModCtrl,
		},
		key.Filter{
			Name:     key.Name("-"),
			Optional: key.ModCtrl,
		},
		key.Filter{
			Name:     key.Name("I"),
			Optional: key.ModCtrl,
		},
		key.Filter{
			Name:     key.Name("R"),
			Optional: key.ModCtrl,
		},
		key.Filter{
			Name: key.Name(key.NameF5),
		},
		key.Filter{
			Name: key.Name(key.NameF6),
		},
	)

	switch event := event.(type) {
	case key.Event:
		if event.State != key.Press {
			return
		}
		if debugHandleGlobalEvents {
			fmt.Println("Event: ", event.Name)
		}
		switch event.Name {
		case "+":
			if event.Modifiers.Contain(key.ModCtrl) {
				if debugHandleGlobalEvents {
					fmt.Printf("Text size increased: %d\n", int(theme.TextSize))
				}
				theme.TextSize += 1
			}
		case "-":
			if event.Modifiers.Contain(key.ModCtrl) {
				if theme.TextSize > 16 {
					if debugHandleGlobalEvents {
						fmt.Printf("Text size decreased: %d\n", int(theme.TextSize))
					}
					theme.TextSize -= 1
				}
			}
		case "I":
			if event.Modifiers.Contain(key.ModCtrl) {
				if debugHandleGlobalEvents {
					fmt.Println("Load program")
				}
				LoadProgramObj(w)
			}
		case "R":
			if event.Modifiers.Contain(key.ModCtrl) {
				if debugHandleGlobalEvents {
					fmt.Println("Reset")
				}
				Reset(w)
			}
		case key.NameF5:
			if debugHandleGlobalEvents {
				fmt.Println("Execute start/stop")
			}
			ExecuteStartStop()
		case key.NameF6:
			if debugHandleGlobalEvents {
				fmt.Println("Execute step")
			}
			ExecuteStep()
		}
	}
}
