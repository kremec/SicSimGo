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
			return material.List(theme, instructionList).Layout(gtx, len(core.Disassembly), func(gtx C, index int) D {
				instruction := core.Disassembly[index]
				instructionAddress := instruction.InstructionAddress.StringHex()
				instructionBytes := fmt.Sprintf("%-6s", strings.ToUpper(hex.EncodeToString(instruction.Instruction.Bytes)))
				instructionOperation := fmt.Sprintf("%-4s", instruction.Instruction.Opcode.String())
				var instructionOperand string
				if instruction.Instruction.Format == core.InstructionFormat2 {
					instructionOperand = fmt.Sprintf("%s,%s", instruction.R1.String(), instruction.R2.String())
				} else if instruction.Instruction.IsJumpInstruction() {
					instructionOperand = instruction.Address.StringHex()
				} else {
					instructionOperand = instruction.Operand.StringHex() + " (" + instruction.Address.StringHex() + ")"
				}
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.End,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return material.Body1(theme, instructionAddress+" : "+instructionBytes+" : "+instructionOperation+" : "+instructionOperand).Layout(gtx)
					}),
				)
			})
		}),
	)
}
