package components

import (
	"fmt"

	"sicsimgo/core/base"
	"sicsimgo/core/loader"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func widthSpacer(width int) layout.FlexChild { // TODO: Also in Dissasembly
	return layout.Rigid(func(gtx C) D {
		return layout.Spacer{Width: unit.Dp(width)}.Layout(gtx)
	})
}
func watchLine(gtx layout.Context, theme *material.Theme, values []string) D {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			value := fmt.Sprintf("%-12s", values[0])
			label := material.Body1(theme, value)
			return label.Layout(gtx)
		}),
		widthSpacer(20),

		layout.Rigid(func(gtx C) D {
			value := fmt.Sprintf("%-12s", values[1])
			label := material.Body1(theme, value)
			return label.Layout(gtx)
		}),
		widthSpacer(20),

		layout.Rigid(func(gtx C) D {
			value := fmt.Sprintf("%-6s", values[2])
			label := material.Body1(theme, value)
			return label.Layout(gtx)
		}),
		widthSpacer(20),

		layout.Rigid(func(gtx C) D {
			value := fmt.Sprintf("%-8s", values[3])
			label := material.Body1(theme, value)
			return label.Layout(gtx)
		}),
	)
}

func Watch(gtx *layout.Context, theme *material.Theme, watchList *widget.List) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(*gtx,
		layout.Rigid(func(gtx C) D {
			return material.H6(theme, "Watch").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return watchLine(gtx, theme, []string{
				"NAME",
				"ADDRESS",
				"DEC",
				"HEX",
			})
		}),

		layout.Flexed(1, func(gtx C) D {
			return material.List(theme, watchList).Layout(gtx, len(loader.SymbolTableList), func(gtx C, index int) D {
				symbol := loader.SymbolTableList[index]
				symbolName := symbol.Name
				symbolAddress := symbol.Address.StringHex()
				var symbolValueDec string
				var symbolValueHex string
				if symbol.DataLength == 1 {
					symbolValue := base.GetByte(symbol.Address)
					symbolValueDec = fmt.Sprintf("%d", int8(symbolValue))
					symbolValueHex = fmt.Sprintf("%02X", symbolValue)
				} else if symbol.DataLength == 3 {
					symbolValue := base.GetWord(symbol.Address)
					symbolValueDec = symbolValue.StringDecSigned()
					symbolValueHex = symbolValue.StringHex()
				}

				return watchLine(gtx, theme, []string{
					symbolName,
					symbolAddress,
					symbolValueDec,
					symbolValueHex,
				})
			})
		}),
	)
}
