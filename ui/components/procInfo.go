package components

import (
	"fmt"
	"sicsimgo/core/proc"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

func ProcInfo(
	gtx *C, theme *material.Theme,
) D {

	var currentInstructionSize int = len(proc.CurrentProcState.Instruction.Bytes)
	var currentInstructionHex string
	for i := 0; i < currentInstructionSize; i++ {
		currentInstructionHex += fmt.Sprintf("%02X ", proc.CurrentProcState.Instruction.Bytes[i])
	}
	var currentInstructionBin string
	for i := 0; i < currentInstructionSize; i++ {
		currentInstructionBin += fmt.Sprintf("%08b ", proc.CurrentProcState.Instruction.Bytes[i])
	}

	var currentInstructionOpcode string = fmt.Sprintf("%02X", proc.CurrentProcState.Instruction.OpcodeByte)

	instructionFormat34 := proc.CurrentProcState.Instruction.Format == proc.InstructionFormat3 || proc.CurrentProcState.Instruction.Format == proc.InstructionFormat4
	var currentBitsNixbpe string

	if instructionFormat34 {
		if proc.CurrentProcState.N {
			currentBitsNixbpe += "n"
		} else {
			currentBitsNixbpe += "-"
		}
		if proc.CurrentProcState.I {
			currentBitsNixbpe += "i"
		} else {
			currentBitsNixbpe += "-"
		}
		if proc.CurrentProcState.X {
			currentBitsNixbpe += "x"
		} else {
			currentBitsNixbpe += "-"
		}
		if proc.CurrentProcState.B {
			currentBitsNixbpe += "b"
		} else {
			currentBitsNixbpe += "-"
		}
		if proc.CurrentProcState.P {
			currentBitsNixbpe += "p"
		} else {
			currentBitsNixbpe += "-"
		}
		if proc.CurrentProcState.E {
			currentBitsNixbpe += "e"
		} else {
			currentBitsNixbpe += "-"
		}
	} else {
		currentBitsNixbpe += "/"
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(*gtx,
		layout.Rigid(func(gtx C) D {
			return material.Body1(theme, "Current instruction bytes:").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return material.Body1(theme, "Hex: "+currentInstructionHex).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return material.Body1(theme, "Bin: "+currentInstructionBin).Layout(gtx)
		}),

		layout.Rigid(layout.Spacer{Height: 10}.Layout),

		layout.Rigid(func(gtx C) D {
			return material.Body1(theme, "Current instruction opcode:").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return material.Body1(theme, "Opcode: "+currentInstructionOpcode).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return material.Body1(theme, "Operation: "+proc.CurrentProcState.Instruction.Opcode.String()).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return material.Body1(theme, "Instruction format: "+proc.CurrentProcState.Instruction.Format.String()).Layout(gtx)
		}),

		layout.Rigid(layout.Spacer{Height: 0}.Layout),

		layout.Rigid(func(gtx C) D {
			if instructionFormat34 {
				return material.Body1(theme, "Bits nixbpe: "+currentBitsNixbpe).Layout(gtx)
			} else {
				return layout.Dimensions{}
			}
		}),
		layout.Rigid(func(gtx C) D {
			if instructionFormat34 {
				return material.Body1(theme, "Address: "+proc.CurrentProcState.Address.StringHex()).Layout(gtx)
			} else {
				return layout.Dimensions{}
			}
		}),
		layout.Rigid(func(gtx C) D {
			if instructionFormat34 && !proc.CurrentProcState.Instruction.Opcode.IsJumpInstruction() {
				return material.Body1(theme, "Operand: "+proc.CurrentProcState.Operand.StringHex()).Layout(gtx)
			} else {
				return layout.Dimensions{}
			}
		}),
	)
}
