package components

import (
	"sicsimgo/core/base"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func DrawRegister(gtx C, theme *material.Theme, name string, value string) D {
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
					return DrawRegister(gtx, theme, "A", base.GetRegisterA().StringHex())
				}),
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "X", base.GetRegisterX().StringHex())
				}),
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "L", base.GetRegisterX().StringHex())
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceAround,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "B", base.GetRegisterB().StringHex())
				}),
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "S", base.GetRegisterS().StringHex())
				}),
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "T", base.GetRegisterT().StringHex())
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceAround,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "F", base.GetRegisterF().StringHex())
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
					return DrawRegister(gtx, theme, "PC", base.GetRegisterPC().StringHex())
				}),
				layout.Rigid(func(gtx C) D {
					return DrawRegister(gtx, theme, "SW", base.GetRegisterSW().StringHex())
				}),
			)
		}),
	)
}
