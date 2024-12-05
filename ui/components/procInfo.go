package components

import (
	"fmt"
	"sicsimgo/core"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

func ProcInfo(
	gtx *C, theme *material.Theme,
) D {

	var currentInstructionSize int = len(core.CurrentProcState.Instruction.Bytes)
	if currentInstructionSize == 0 {
		return layout.Dimensions{}
	}

	var currentInstructionHex string
	for i := 0; i < currentInstructionSize; i++ {
		currentInstructionHex += fmt.Sprintf("%02X ", core.CurrentProcState.Instruction.Bytes[i])
	}
	var currentInstructionBin string
	for i := 0; i < currentInstructionSize; i++ {
		currentInstructionBin += fmt.Sprintf("%08b ", core.CurrentProcState.Instruction.Bytes[i])
	}

	var currentInstructionOpcode string = fmt.Sprintf("%02X", core.CurrentProcState.Instruction.Bytes[0])

	instructionFormat34 := core.CurrentProcState.Instruction.Format == core.InstructionFormat3 || core.CurrentProcState.Instruction.Format == core.InstructionFormat4
	var currentBitsNixbpe string

	if instructionFormat34 {
		if core.CurrentProcState.N {
			currentBitsNixbpe += "n"
		} else {
			currentBitsNixbpe += "-"
		}
		if core.CurrentProcState.I {
			currentBitsNixbpe += "i"
		} else {
			currentBitsNixbpe += "-"
		}
		if core.CurrentProcState.X {
			currentBitsNixbpe += "x"
		} else {
			currentBitsNixbpe += "-"
		}
		if core.CurrentProcState.B {
			currentBitsNixbpe += "b"
		} else {
			currentBitsNixbpe += "-"
		}
		if core.CurrentProcState.P {
			currentBitsNixbpe += "p"
		} else {
			currentBitsNixbpe += "-"
		}
		if core.CurrentProcState.E {
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
			return material.Body1(theme, "Operation: "+core.CurrentProcState.Instruction.Opcode.String()).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return material.Body1(theme, "Instruction format: "+core.CurrentProcState.Instruction.Format.String()).Layout(gtx)
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
				return material.Body1(theme, "Address: "+core.CurrentProcState.Instruction.Address.StringHex()).Layout(gtx)
			} else {
				return layout.Dimensions{}
			}
		}),
		layout.Rigid(func(gtx C) D {
			if instructionFormat34 && !core.CurrentProcState.Instruction.IsJumpInstruction() {
				return material.Body1(theme, "Operand: "+core.CurrentProcState.Instruction.Operand.StringHex()).Layout(gtx)
			} else {
				return layout.Dimensions{}
			}
		}),
	)
}
