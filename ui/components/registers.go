package components

import (
	"sicsimgo/core"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func DrawRegister(gtx layout.Context, theme *material.Theme, name string, value string) layout.Dimensions {
	registerMargins := layout.Inset{
		Top:    unit.Dp(0),
		Bottom: unit.Dp(0),
		Right:  unit.Dp(10),
		Left:   unit.Dp(0),
	}

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return registerMargins.Layout(gtx, func(gtx C) D {
				return material.Body1(theme, name).Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return material.Body1(theme, value).Layout(gtx)
		}),
	)
}

func Registers(gtx C, theme *material.Theme) D {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return material.H6(theme, "Registers").Layout(gtx)
		}),

		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceAround,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "A", core.GetRegisterA().StringHex())
				}),
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "X", core.GetRegisterX().StringHex())
				}),
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "L", core.GetRegisterX().StringHex())
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceAround,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "B", core.GetRegisterB().StringHex())
				}),
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "S", core.GetRegisterS().StringHex())
				}),
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "T", core.GetRegisterT().StringHex())
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceAround,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "F", core.GetRegisterF().StringHex())
				}),
			)
		}),

		layout.Rigid(layout.Spacer{Height: 10}.Layout),

		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceAround,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "PC", core.GetRegisterPC().StringHex())
				}),
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "SW", core.GetRegisterSW().StringHex())
				}),
			)
		}),
	)
}
