package components

import (
	"fmt"
	"sicsimgo/core"
	"sicsimgo/core/proc"

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

	instructionFormat34 := core.CurrentProcState.Instruction.Format == proc.InstructionFormat3 || core.CurrentProcState.Instruction.Format == proc.InstructionFormat4
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
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(*gtx,
		layout.Rigid(func(gtx C) D {
			return material.H6(theme, "Next instruction").Layout(gtx)
		}),

		layout.Rigid(func(gtx C) D {
			return material.Body1(theme, fmt.Sprintf("Opcode (Operation): %s (%s)", currentInstructionOpcode, core.CurrentProcState.Instruction.Opcode.String())).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return material.Body1(theme, "Format: "+core.CurrentProcState.Instruction.Format.String()).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			if instructionFormat34 {
				return material.Body1(theme, "nixbpe: "+currentBitsNixbpe).Layout(gtx)
			} else {
				return layout.Dimensions{}
			}
		}),
		layout.Rigid(func(gtx C) D {
			return material.Body1(theme, "Hex: "+currentInstructionHex).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return material.Body1(theme, "Bin: "+currentInstructionBin).Layout(gtx)
		}),
	)
}
