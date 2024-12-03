package loader

import (
	"sicsimgo/core"
	"sicsimgo/core/proc"
	"sicsimgo/core/units"
)

type Instruction struct {
	InstructionAddress units.Int24
	Instruction        proc.Instruction
	Address            units.Int24
	Operand            units.Int24
	R1, R2             core.RegisterId
}

var Instructions []Instruction

func GetInstructionsFromTextRecord(codeAddress units.Int24, code []byte) []Instruction {
	instructions := []Instruction{}

	byteIndex := 0
	relativeAddress := units.Int24{}
	for byteIndex < len(code) {
		instruction := Instruction{}

		byte1 := code[byteIndex]
		opcodeByte := byte1 & 0xFC
		instruction.Instruction.OpcodeByte = opcodeByte
		opcode := proc.Opcode(opcodeByte)
		instruction.Instruction.Opcode = opcode

		instructionFormatFromOpcode, err := proc.GetInstructionFormatFromOpcode(opcode, byte1)
		if err != nil {
			// Invalid opcode
			panic(err)
		}
		instruction.Instruction.Format = instructionFormatFromOpcode

		instructionBytes := []byte{byte1}
		if instructionFormatFromOpcode == proc.InstructionFormat2 {
			byte2 := code[byteIndex+1]

			r1, r2 := proc.GetR1R2FromByte(byte2)
			instruction.R1 = r1
			instruction.R2 = r2

			instructionBytes = append(instructionBytes, byte2)
		} else if instructionFormatFromOpcode == proc.InstructionFormatSIC {
			byte2 := code[byteIndex+1]
			byte3 := code[byteIndex+2]
			instructionBytes = append(instructionBytes, byte2)
			instructionBytes = append(instructionBytes, byte3)
		} else if instructionFormatFromOpcode == proc.InstructionFormat34 {
			byte2 := code[byteIndex+1]
			byte3 := code[byteIndex+2]
			instructionBytes = append(instructionBytes, byte2)
			instructionBytes = append(instructionBytes, byte3)

			instructionType := proc.GetInstructionFormat3Or4FromByte(byte2)
			if instructionType == proc.InstructionFormat3 {
				instruction.Instruction.Format = proc.InstructionFormat3
			} else if instructionType == proc.InstructionFormat4 {
				byte4 := code[byteIndex+3]
				instructionBytes = append(instructionBytes, byte4)

				instruction.Instruction.Format = proc.InstructionFormat4
			}
		}
		instruction.Instruction.Bytes = instructionBytes
		instruction.InstructionAddress = codeAddress.Add(relativeAddress)

		byteIndex += len(instructionBytes)
		for i := 0; i < len(instructionBytes); i++ {
			relativeAddress = relativeAddress.Add(units.Int24{0x00, 0x00, 0x01})
		}

		nextInstructionAddress := codeAddress.Add(relativeAddress)
		if instructionFormatFromOpcode == proc.InstructionFormat34 {
			operand, address := instruction.Instruction.GetOperandAddress(nextInstructionAddress)
			instruction.Address = address
			instruction.Operand = operand
		}

		instructions = append(instructions, instruction)
	}

	return instructions
}
