package core

import (
	"sicsimgo/core/base"
	"sicsimgo/core/units"
)

/*
DEFINITIONS
*/
type DisassemblyInstruction struct {
	InstructionAddress units.Int24
	Instruction        Instruction
	Address            units.Int24
	Operand            units.Int24
	R1, R2             base.RegisterId
}

/*
IMPLEMENTATION
*/
var Disassembly []DisassemblyInstruction

/*
OPERATIONS
*/
func GetDisassemblyInstructionsFromTextRecord(codeAddress units.Int24, code []byte) []DisassemblyInstruction {
	instructions := []DisassemblyInstruction{}

	byteIndex := 0
	relativeAddress := units.Int24{}
	for byteIndex < len(code) {
		instruction := DisassemblyInstruction{}

		byte1 := code[byteIndex]
		opcode := GetOpcode(byte1)
		instruction.Instruction.Opcode = opcode

		instructionFormatFromOpcode, err := GetInstructionFormat(byte1)
		if err != nil {
			// Invalid opcode
			panic(err)
		}
		instruction.Instruction.Format = instructionFormatFromOpcode

		instructionBytes := []byte{byte1}
		if instructionFormatFromOpcode == InstructionFormat2 {
			byte2 := code[byteIndex+1]

			r1, r2 := GetR1R2FromByte(byte2)
			instruction.R1 = r1
			instruction.R2 = r2

			instructionBytes = append(instructionBytes, byte2)
		} else if instructionFormatFromOpcode == InstructionFormatSIC {
			byte2 := code[byteIndex+1]
			byte3 := code[byteIndex+2]
			instructionBytes = append(instructionBytes, byte2)
			instructionBytes = append(instructionBytes, byte3)
		} else if instructionFormatFromOpcode == InstructionFormat34 {
			byte2 := code[byteIndex+1]
			byte3 := code[byteIndex+2]
			instructionBytes = append(instructionBytes, byte2)
			instructionBytes = append(instructionBytes, byte3)

			instructionType := GetInstructionFormat34(byte2)
			if instructionType == InstructionFormat3 {
				instruction.Instruction.Format = InstructionFormat3
			} else if instructionType == InstructionFormat4 {
				byte4 := code[byteIndex+3]
				instructionBytes = append(instructionBytes, byte4)

				instruction.Instruction.Format = InstructionFormat4
			}
		}
		instruction.Instruction.Bytes = instructionBytes
		instruction.InstructionAddress = codeAddress.Add(relativeAddress)

		byteIndex += len(instructionBytes)
		for i := 0; i < len(instructionBytes); i++ {
			relativeAddress = relativeAddress.Add(units.Int24{0x00, 0x00, 0x01})
		}

		nextInstructionAddress := codeAddress.Add(relativeAddress)
		if instructionFormatFromOpcode == InstructionFormat34 {
			operand, address := instruction.Instruction.GetOperandAddress(nextInstructionAddress)
			instruction.Address = address
			instruction.Operand = operand
		}

		instructions = append(instructions, instruction)
	}

	return instructions
}
