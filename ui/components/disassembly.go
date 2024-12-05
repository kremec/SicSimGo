package components

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"strings"

	"sicsimgo/core"
	"sicsimgo/core/base"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/image/colornames"
)

func WidthSpacer(gtx layout.Context, width int) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		return layout.Spacer{Width: unit.Dp(width)}.Layout(gtx)
	})
}
func InstructionLine(gtx layout.Context, theme *material.Theme, values []string, selected bool) D {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			value := fmt.Sprintf("%-8s", values[0])
			label := material.Body1(theme, value)
			if selected {
				label.Color = color.NRGBA(colornames.Red)
			}
			return label.Layout(gtx)
		}),
		WidthSpacer(gtx, 20),

		layout.Rigid(func(gtx C) D {
			value := fmt.Sprintf("%-8s", values[1])
			label := material.Body1(theme, value)
			if selected {
				label.Color = color.NRGBA(colornames.Red)
			}
			return label.Layout(gtx)
		}),
		WidthSpacer(gtx, 20),

		layout.Rigid(func(gtx C) D {
			value := fmt.Sprintf("%-9s", values[2])
			label := material.Body1(theme, value)
			if selected {
				label.Color = color.NRGBA(colornames.Red)
			}
			return label.Layout(gtx)
		}),
		WidthSpacer(gtx, 20),

		layout.Rigid(func(gtx C) D {
			value := fmt.Sprintf("%s", values[3])
			label := material.Body1(theme, value)
			if selected {
				if core.CurrentProcState.Instruction.IsFormatSIC34() && core.CurrentProcState.Instruction.AbsoluteAddressingMode != core.ImmediateAbsoluteAddressing {
					label.Color = color.NRGBA(colornames.Darkorchid)
				} else {
					label.Color = color.NRGBA(colornames.Red)
				}
			}
			return label.Layout(gtx)
		}),
	)
}

func Disassembly(gtx *layout.Context, theme *material.Theme, instructionList *widget.List) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(*gtx,
		layout.Rigid(func(gtx C) D {
			return material.H6(theme, "Disassembly").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return InstructionLine(gtx, theme, []string{
				"ADDRESS",
				"BYTES",
				"OPERATION",
				"OPERAND (ADRRESS)",
			}, false)
		}),

		layout.Flexed(1, func(gtx C) D {
			return material.List(theme, instructionList).Layout(gtx, len(core.InstructionList), func(gtx C, index int) D {
				instruction := core.InstructionList[index]
				instructionAddress := instruction.InstructionAddress.StringHex()
				instructionBytes := fmt.Sprintf("%-8s", strings.ToUpper(hex.EncodeToString(instruction.Bytes)))
				var instructionOperation string
				if instruction.Directive == core.DirectiveBYTE {
					instructionOperation = fmt.Sprintf("%-4s", core.DirectiveBYTE)
				} else {
					if instruction.Format == core.InstructionFormat4 {
						instructionOperation = "+"
					} else {
						instructionOperation = " "
					}
					instructionOperation += fmt.Sprintf("%-4s", instruction.Opcode.String())
				}
				var instructionOperand string
				if instruction.Directive == core.DirectiveBYTE {
					instructionOperand = ""
				} else if instruction.Format == core.InstructionFormat2 {
					instructionOperand = fmt.Sprintf("%s,%s", instruction.R1.String(), instruction.R2.String())
				} else if instruction.IsJumpInstruction() || instruction.IsStoreInstruction() {
					instructionOperand = instruction.Address.StringHex()
				} else {
					instructionOperand = instruction.Operand.StringHex() + " (" + instruction.Address.StringHex() + ")"
				}

				instructionSelected := instruction.InstructionAddress.Compare(base.GetRegisterPC()) == 0
				return InstructionLine(gtx, theme, []string{
					instructionAddress,
					instructionBytes,
					instructionOperation,
					instructionOperand,
				}, instructionSelected)
			})
		}),
	)
}
