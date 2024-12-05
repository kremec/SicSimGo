package components

import (
	"encoding/hex"
	"fmt"
	"strings"

	"sicsimgo/core"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func Disassembly(gtx *layout.Context, theme *material.Theme, instructionList *widget.List) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(*gtx,
		layout.Rigid(func(gtx C) D {
			return material.H6(theme, "Disassembly").Layout(gtx)
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
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.End,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return material.Body1(theme, instructionAddress+" : "+instructionBytes+" → "+instructionOperation+" : "+instructionOperand).Layout(gtx)
					}),
				)
			})
		}),
	)
}
