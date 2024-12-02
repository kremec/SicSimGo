package components

import (
	"sicsimgo/core"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func Memory(gtx *layout.Context, theme *material.Theme, memoryList *widget.List) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(*gtx,
		layout.Rigid(func(gtx C) D {
			return material.H6(theme, "Memory").Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return material.List(theme, memoryList).Layout(gtx, 65536, func(gtx C, index int) D {
				address := core.ToAddress(uint32(index) * 16)
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.End,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return material.Body1(theme, core.StringAddress(address)+"  :  ").Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return material.Body1(theme, core.String16Bytes(address)).Layout(gtx)
					}),
				)
			})
		}),
	)
}
