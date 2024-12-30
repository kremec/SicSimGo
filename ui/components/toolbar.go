package components

import (
	"sicsimgo/core"
	"sicsimgo/core/loader"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func toolbarButton(theme *material.Theme, button *widget.Clickable, label string) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		return layout.Inset{
			Top:    unit.Dp(0),
			Bottom: unit.Dp(0),
			Right:  unit.Dp(1),
			Left:   unit.Dp(1),
		}.Layout(gtx, func(gtx C) D {
			return material.Button(theme, button, label).Layout(gtx)
		})
	})
}

func Toolbar(gtx C, theme *material.Theme, LoadProgramButton, ExecuteStepButton, ExecuteStartButton, ResetSimButton, OutputObjFileButton, OutputLstFileButton *widget.Clickable) D {

	ExecuteState := func() string {
		if core.SimExecuteState == core.ExecuteStartState {
			return "STOP"
		} else {
			return "START"
		}
	}()

	return layout.Inset{
		Top:    unit.Dp(2),
		Bottom: unit.Dp(2),
		Right:  unit.Dp(2),
		Left:   unit.Dp(2),
	}.Layout(gtx, func(gtx C) D {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			toolbarButton(theme, LoadProgramButton, "LOAD"),
			toolbarButton(theme, ResetSimButton, "RESET"),
			toolbarButton(theme, ExecuteStepButton, "STEP"),
			toolbarButton(theme, ExecuteStartButton, ExecuteState),
			layout.Flexed(1, func(gtx C) D {
				return layout.Spacer{}.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				if core.LoadedProgramTypeState == loader.Assembly {
					return layout.Flex{}.Layout(gtx,
						toolbarButton(theme, OutputLstFileButton, "LST"),
					)
				}
				return D{}
			}),
			layout.Rigid(func(gtx C) D {
				if core.LoadedProgramTypeState == loader.Assembly {
					return layout.Flex{}.Layout(gtx,
						toolbarButton(theme, OutputObjFileButton, "OBJ"),
					)
				}
				return D{}
			}),
		)
	})

}
