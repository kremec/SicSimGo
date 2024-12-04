package core

import (
	"encoding/hex"
	"fmt"
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
DEBUG
*/
const debugGetDisassemblyInstructionsFromTextRecord bool = false

/*
OPERATIONS
*/
func GetDisassemblyInstructionsFromTextRecord(codeAddress units.Int24, code []byte) ([]DisassemblyInstruction, []byte) {
	instructions := []DisassemblyInstruction{}

	byteIndex := 0
	relativeAddress := units.Int24{}
	for byteIndex < len(code) {
		instruction := DisassemblyInstruction{}

		if debugGetDisassemblyInstructionsFromTextRecord {
			fmt.Printf("Looking at byte index %d\n", byteIndex)
		}

		byte1 := code[byteIndex]
		opcode := GetOpcode(byte1)
		instruction.Instruction.Opcode = opcode

		instructionFormatFromOpcode, err := GetInstructionFormat(byte1)
		if err != nil {
			// Invalid opcode
			panic(err)
		}
		instruction.Instruction.Format = instructionFormatFromOpcode

		if debugGetDisassemblyInstructionsFromTextRecord {
			fmt.Printf("    Opcode: %s\n", opcode.String())
			fmt.Printf("    Instruction format: %s\n", instructionFormatFromOpcode.String())
		}

		instructionBytes := []byte{byte1}
		if instructionFormatFromOpcode == InstructionFormat2 {
			// Check for incomplete instruction
			if (byteIndex + 1) >= len(code) {
				return instructions, []byte{byte1}
			}
			byte2 := code[byteIndex+1]

			r1, r2 := GetR1R2FromByte(byte2)
			instruction.R1 = r1
			instruction.R2 = r2

			instructionBytes = append(instructionBytes, byte2)
		} else if instructionFormatFromOpcode == InstructionFormatSIC {
			// Check for incomplete instruction
			if (byteIndex + 1) >= len(code) {
				return instructions, []byte{byte1}
			}
			byte2 := code[byteIndex+1]
			if (byteIndex + 2) >= len(code) {
				return instructions, []byte{byte1, byte2}
			}
			byte3 := code[byteIndex+2]

			instructionBytes = append(instructionBytes, byte2)
			instructionBytes = append(instructionBytes, byte3)
		} else if instructionFormatFromOpcode == InstructionFormat34 {
			// Check for incomplete instruction
			if (byteIndex + 1) >= len(code) {
				return instructions, []byte{byte1}
			}
			byte2 := code[byteIndex+1]
			if (byteIndex + 2) >= len(code) {
				return instructions, []byte{byte1, byte2}
			}
			byte3 := code[byteIndex+2]

			instructionBytes = append(instructionBytes, byte2)
			instructionBytes = append(instructionBytes, byte3)

			instructionType := GetInstructionFormat34(byte2)
			if instructionType == InstructionFormat3 {
				instruction.Instruction.Format = InstructionFormat3
			} else if instructionType == InstructionFormat4 {
				// Check for incomplete instruction
				if (byteIndex + 3) >= len(code) {
					return instructions, []byte{byte1}
				}
				byte4 := code[byteIndex+3]

				instructionBytes = append(instructionBytes, byte4)
				instruction.Instruction.Format = InstructionFormat4
			}
		}
		instruction.Instruction.Bytes = instructionBytes
		instruction.InstructionAddress = codeAddress.Add(relativeAddress)

		if debugGetDisassemblyInstructionsFromTextRecord {
			fmt.Printf("    Instruction bytes: %s\n", hex.EncodeToString(instructionBytes))
			fmt.Printf("    Instruction address: %s\n", instruction.InstructionAddress.StringHex())
		}

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

	return instructions, []byte{}
}
