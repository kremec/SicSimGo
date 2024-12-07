package components

import (
	"fmt"
	"image/color"
	"sicsimgo/core"
	"sicsimgo/core/base"
	"sicsimgo/core/proc"
	"sicsimgo/core/units"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/image/colornames"
)

func MemoryByte(gtx layout.Context, theme *material.Theme, address units.Int24, value byte, instructionAddressSelected bool, operandAddressSelected bool) D {
	label := material.Body1(theme, fmt.Sprintf("%02X", value))
	if instructionAddressSelected {
		label.Color = color.NRGBA(colornames.Red)
	} else if operandAddressSelected {
		label.Color = color.NRGBA(colornames.Darkorchid)
	}
	return label.Layout(gtx)
}
func MemoryLine(gtx layout.Context, theme *material.Theme, address units.Int24, values []byte, instructionAddressSelection []bool, operandAddressSelection []bool) D {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			value := base.StringAddress(address)
			return material.Body1(theme, value).Layout(gtx)
		}),
		WidthSpacer(gtx, 20),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[0], instructionAddressSelection[0], operandAddressSelection[0])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[1], instructionAddressSelection[1], operandAddressSelection[1])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[2], instructionAddressSelection[2], operandAddressSelection[2])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[3], instructionAddressSelection[3], operandAddressSelection[3])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[4], instructionAddressSelection[4], operandAddressSelection[4])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[5], instructionAddressSelection[5], operandAddressSelection[5])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[6], instructionAddressSelection[6], operandAddressSelection[6])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[7], instructionAddressSelection[7], operandAddressSelection[7])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[8], instructionAddressSelection[8], operandAddressSelection[8])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[9], instructionAddressSelection[9], operandAddressSelection[9])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[10], instructionAddressSelection[10], operandAddressSelection[10])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[11], instructionAddressSelection[11], operandAddressSelection[11])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[12], instructionAddressSelection[12], operandAddressSelection[12])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[13], instructionAddressSelection[13], operandAddressSelection[13])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[14], instructionAddressSelection[14], operandAddressSelection[14])
		}),
		WidthSpacer(gtx, 10),

		layout.Rigid(func(gtx C) D {
			return MemoryByte(gtx, theme, address, values[15], instructionAddressSelection[15], operandAddressSelection[15])
		}),
		WidthSpacer(gtx, 10),
	)
}

func Memory(gtx *layout.Context, theme *material.Theme, memoryList *widget.List) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(*gtx,
		layout.Rigid(func(gtx C) D {
			return material.H6(theme, "Memory").Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {

			// PC-instruction selection adresses
			instructionAddresses := []units.Int24{}
			pcAddress := core.CurrentProcState.Instruction.InstructionAddress
			j := units.Int24{}
			for i := 0; i < len(core.CurrentProcState.Instruction.Bytes); i++ {
				instructionAddresses = append(instructionAddresses, pcAddress.Add(j))
				j = j.Add(units.Int24{0x00, 0x00, 0x01})
			}

			// Operand address selection adresses
			operandAddresses := []units.Int24{}
			if core.CurrentProcState.Instruction.IsFormatSIC34() && core.CurrentProcState.Instruction.AbsoluteAddressingMode != proc.ImmediateAbsoluteAddressing {
				operandAddress := core.CurrentProcState.Instruction.Address
				j = units.Int24{}
				for i := 0; i < 3; i++ {
					operandAddresses = append(operandAddresses, operandAddress.Add(j))
					j = j.Add(units.Int24{0x00, 0x00, 0x01})
				}
			}

			return material.List(theme, memoryList).Layout(gtx, 65536, func(gtx C, index int) D {
				address := base.ToAddress(uint32(index) * 16)
				memoryLineAddress := address.ToUint32()

				// PC-instruction selection
				instructionAddressSelection := make([]bool, 16, 16)
				for _, instructionAddress := range instructionAddresses {
					if memoryLineAddress <= instructionAddress.ToUint32() && instructionAddress.ToUint32() <= memoryLineAddress+15 {
						position := int(instructionAddress.ToUint32() - memoryLineAddress)
						instructionAddressSelection[position] = true
					}
				}

				// Operand address selection
				operandAddressSelection := make([]bool, 16, 16)
				for _, operandAddress := range operandAddresses {
					if memoryLineAddress <= operandAddress.ToUint32() && operandAddress.ToUint32() <= memoryLineAddress+15 {
						position := int(operandAddress.ToUint32() - memoryLineAddress)
						operandAddressSelection[position] = true
					}
				}

				return MemoryLine(gtx, theme, address, base.GetSlice16(address), instructionAddressSelection, operandAddressSelection)
			})
		}),
	)
}
