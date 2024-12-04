package ui

import (
	_ "embed"
	"sicsimgo/core"
	"sicsimgo/internal"
	"sicsimgo/ui/components"
	"sicsimgo/ui/events"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

//go:embed FiraMono-Regular.ttf
var fontBytes []byte

func loadFont(theme *material.Theme) error {

	faces, err := opentype.ParseCollection(fontBytes)
	if err != nil {
		return err
	}

	collection := gofont.Collection()
	theme.Shaper = text.NewShaper(text.WithCollection(append(collection, faces...)))
	theme.Face = "Fira Mono"

	return nil
}

func DrawWindow(w *app.Window) error {
	var ops op.Ops

	theme := material.NewTheme()
	loadFont(theme)

	var LoadProgramButton widget.Clickable
	var ExecuteStepButton widget.Clickable
	var ExecuteStartStopButton widget.Clickable
	var ResetSimButton widget.Clickable

	memoryList := widget.List{
		List: layout.List{Axis: layout.Vertical},
	}
	instructionList := widget.List{
		List: layout.List{Axis: layout.Vertical},
	}

	core.InitProcState()

	for {
		switch e := w.Event().(type) {

		// Application rerender
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			events.HandleGlobalEvents(gtx, theme)

			if LoadProgramButton.Clicked(gtx) {
				core.ResetSim()
				core.InitProcState()

				go func() {
					programName := core.OpenObjectFile()
					internal.SetWindowTitle(programName, w)
					core.InitProcState()
				}()
			}
			if ExecuteStepButton.Clicked(gtx) {
				go core.ExecuteNextInstruction()
			}
			if ExecuteStartStopButton.Clicked(gtx) {
				core.SimExecuteState = !core.SimExecuteState
				go func() {
					for core.SimExecuteState == core.ExecuteStartState {
						core.ExecuteNextInstruction()
					}
				}()
			}
			if ResetSimButton.Clicked(gtx) {
				internal.ResetWindowTitle(w)
				go func() {
					core.ResetSim()
					core.InitProcState()
				}()
			}

			layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return components.Toolbar(gtx, theme, &LoadProgramButton, &ExecuteStepButton, &ExecuteStartStopButton, &ResetSimButton)
				}),

				layout.Flexed(1, func(gtx C) D {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Flexed(0.4, func(gtx C) D {

							return layout.Flex{
								Axis:      layout.Vertical,
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Flex{
										Axis:      layout.Vertical,
										Alignment: layout.Middle,
									}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return layout.Inset{
												Top:    unit.Dp(5),
												Bottom: unit.Dp(0),
												Right:  unit.Dp(0),
												Left:   unit.Dp(5),
											}.Layout(gtx, func(gtx C) D {
												return components.Registers(gtx, theme)
											})
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{
												Top:    unit.Dp(5),
												Bottom: unit.Dp(0),
												Right:  unit.Dp(5),
												Left:   unit.Dp(0),
											}.Layout(gtx, func(gtx C) D {
												return components.ProcInfo(
													&gtx, theme,
												)
											})
										}),
									)
								}),
								layout.Rigid(func(gtx C) D {
									return component.Divider(theme).Layout(gtx)
								}),
								layout.Flexed(1, func(gtx C) D {
									return layout.Inset{
										Top:    unit.Dp(0),
										Bottom: unit.Dp(5),
										Right:  unit.Dp(5),
										Left:   unit.Dp(5),
									}.Layout(gtx, func(gtx C) D {
										return components.Disassembly(&gtx, theme, &instructionList)
									})
								}),
							)
						}),

						layout.Flexed(0.6, func(gtx C) D {
							return layout.Flex{
								Axis:      layout.Vertical,
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Flexed(0.3, func(gtx C) D {
									return layout.Inset{
										Top:    unit.Dp(5),
										Bottom: unit.Dp(0),
										Right:  unit.Dp(5),
										Left:   unit.Dp(5),
									}.Layout(gtx, func(gtx C) D {
										return material.H6(theme, "Watch").Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return component.Divider(theme).Layout(gtx)
								}),
								layout.Flexed(0.7, func(gtx C) D {
									return layout.Inset{
										Top:    unit.Dp(0),
										Bottom: unit.Dp(5),
										Right:  unit.Dp(5),
										Left:   unit.Dp(5),
									}.Layout(gtx, func(gtx C) D {
										return components.Memory(&gtx, theme, &memoryList)
									})
								}),
							)
						}),
					)
				}),
			)

			e.Frame(gtx.Ops)

		// Application exit
		case app.DestroyEvent:
			return e.Err
		}
	}
}
