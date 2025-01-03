package ui

import (
	_ "embed"
	"sicsimgo/core"
	"sicsimgo/core/base"
	"sicsimgo/core/loader"
	"sicsimgo/internal"
	"sicsimgo/ui/components"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
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

func OpenProgramFile(w *app.Window) {
	core.ResetSim()

	go func() {
		programName, startPC, loadedProgramType := loader.OpenAsmObjFile()
		if loadedProgramType == loader.None {
			core.ResetSim()
			internal.ResetWindowTitle(w)
			return
		}
		base.SetRegisterPC(startPC)
		core.LoadedProgramTypeState = loadedProgramType
		core.UpdateProcState(base.GetRegisterPC())
		internal.SetWindowTitle(programName, w)
	}()
}
func ExecuteStep() {
	go core.ExecuteNextInstruction()
}
func ExecuteStartStop() {
	core.SimExecuteState = !core.SimExecuteState
	go func() {
		for core.SimExecuteState == core.ExecuteStartState {
			core.ExecuteNextInstruction()
		}
	}()
}
func Reset(w *app.Window) {
	internal.ResetWindowTitle(w)
	go func() {
		core.ResetSim()
	}()
}
func OutputLstFile() {
	go func() {
		loader.OutputLstFile()
	}()
}
func OutputObjFile() {
	go func() {
		loader.OutputObjFile()
	}()
}

func DrawWindow(w *app.Window) error {
	var ops op.Ops

	theme := material.NewTheme()
	loadFont(theme)

	var LoadProgramButton widget.Clickable
	var ExecuteStepButton widget.Clickable
	var ExecuteStartStopButton widget.Clickable
	var ResetSimButton widget.Clickable
	var OutputObjFileButton widget.Clickable
	var OutputLstFileButton widget.Clickable

	memoryList := widget.List{
		List: layout.List{Axis: layout.Vertical},
	}
	instructionList := widget.List{
		List: layout.List{Axis: layout.Vertical},
	}
	watchList := widget.List{
		List: layout.List{Axis: layout.Vertical},
	}

	mainSplit := Split{
		Ratio: -0.2,
	}
	vSplitLeft := Split{
		Ratio: -0.2,
	}
	vSplitRight := Split{
		Ratio: -0.7,
	}

	for {
		switch e := w.Event().(type) {

		// Application rerender
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			HandleGlobalEvents(gtx, theme, w)

			if LoadProgramButton.Clicked(gtx) {
				OpenProgramFile(w)
			}
			if ExecuteStepButton.Clicked(gtx) {
				ExecuteStep()
			}
			if ExecuteStartStopButton.Clicked(gtx) {
				ExecuteStartStop()
			}
			if ResetSimButton.Clicked(gtx) {
				Reset(w)
			}
			if OutputLstFileButton.Clicked(gtx) {
				OutputLstFile()
			}
			if OutputObjFileButton.Clicked(gtx) {
				OutputObjFile()
			}

			layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return components.Toolbar(gtx, theme, &LoadProgramButton, &ExecuteStepButton, &ExecuteStartStopButton, &ResetSimButton, &OutputObjFileButton, &OutputLstFileButton)
				}),

				layout.Flexed(1, func(gtx C) D {
					return mainSplit.HLayout(gtx,
						func(gtx C) D {
							return vSplitLeft.VLayout(gtx,
								func(gtx C) D {
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
												Left:   unit.Dp(5),
											}.Layout(gtx, func(gtx C) D {
												return components.ProcInfo(
													&gtx, theme,
												)
											})
										}),
									)
								},
								func(gtx C) D {
									return layout.Inset{
										Top:    unit.Dp(0),
										Bottom: unit.Dp(5),
										Right:  unit.Dp(5),
										Left:   unit.Dp(5),
									}.Layout(gtx, func(gtx C) D {
										return components.Disassembly(&gtx, theme, &instructionList)
									})
								},
							)
						},
						func(gtx C) D {
							return vSplitRight.VLayout(gtx,
								func(gtx C) D {
									return layout.Inset{
										Top:    unit.Dp(5),
										Bottom: unit.Dp(0),
										Right:  unit.Dp(5),
										Left:   unit.Dp(5),
									}.Layout(gtx, func(gtx C) D {
										return components.Watch(&gtx, theme, &watchList)
									})
								},
								func(gtx C) D {
									return layout.Inset{
										Top:    unit.Dp(0),
										Bottom: unit.Dp(5),
										Right:  unit.Dp(5),
										Left:   unit.Dp(5),
									}.Layout(gtx, func(gtx C) D {
										return components.Memory(&gtx, theme, &memoryList)
									})
								},
							)
						},
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
