package bytecode

import (
	"fmt"
	"sicsimgo/core/proc"
	"sicsimgo/core/units"
)

/*
DEBUG
*/
const debugGetDisassemblyInstructionsFromTextRecord bool = false

/*
OPERATIONS
*/
func GetInstructionsFromBinary(codeAddress units.Int24, binaryCode []byte) (map[units.Int24]proc.Instruction, []byte) {
	disassemblyInstructions := make(map[units.Int24]proc.Instruction)

	byteIndex := 0
	relativeAddress := units.Int24{}
	for byteIndex < len(binaryCode) {
		instruction := proc.Instruction{}

		if debugGetDisassemblyInstructionsFromTextRecord {
			fmt.Printf("Looking at byte index %d\n", byteIndex)
		}

		byte1 := binaryCode[byteIndex]
		opcode := proc.GetOpcode(byte1)
		instruction.Opcode = opcode

		instructionFormatFromOpcode, err := proc.GetInstructionFormat(byte1)
		if err != nil {
			// Invalid opcode, perhaps data between instruction bytes
			if debugGetDisassemblyInstructionsFromTextRecord {
				fmt.Printf("    INVALID OPCODE: %s\n", opcode.String())
			}

			// Add byte as data to instructions and continue to next byte in code
			instruction := proc.Instruction{
				InstructionAddress: codeAddress.Add(relativeAddress),
				Format:             proc.InstructionUnknown,
				Bytes:              []byte{byte1},
				Directive:          proc.DirectiveBYTE,
			}

			disassemblyInstructions[instruction.InstructionAddress] = instruction
			byteIndex += 1
			relativeAddress = relativeAddress.Add(units.Int24{0x00, 0x00, 0x01})
			continue
		}
		instruction.Format = instructionFormatFromOpcode

		if debugGetDisassemblyInstructionsFromTextRecord {
			fmt.Printf("    Opcode: %s\n", opcode.String())
			fmt.Printf("    Instruction format: %s\n", instructionFormatFromOpcode.String())
		}

		instructionBytes := []byte{byte1}
		if instructionFormatFromOpcode == proc.InstructionFormat2 {
			// Check for incomplete instruction
			if (byteIndex + 1) >= len(binaryCode) {
				return disassemblyInstructions, []byte{byte1}
			}
			byte2 := binaryCode[byteIndex+1]

			r1, r2 := proc.GetR1R2FromByte(byte2)
			instruction.R1 = r1
			instruction.R2 = r2

			instructionBytes = append(instructionBytes, byte2)
		} else if instructionFormatFromOpcode == proc.InstructionFormatSIC {
			// Check for incomplete instruction
			if (byteIndex + 1) >= len(binaryCode) {
				return disassemblyInstructions, []byte{byte1}
			}
			byte2 := binaryCode[byteIndex+1]
			if (byteIndex + 2) >= len(binaryCode) {
				return disassemblyInstructions, []byte{byte1, byte2}
			}
			byte3 := binaryCode[byteIndex+2]

			instructionBytes = append(instructionBytes, byte2)
			instructionBytes = append(instructionBytes, byte3)
		} else if instructionFormatFromOpcode == proc.InstructionFormat34 {
			// Check for incomplete instruction
			if (byteIndex + 1) >= len(binaryCode) {
				return disassemblyInstructions, []byte{byte1}
			}
			byte2 := binaryCode[byteIndex+1]
			if (byteIndex + 2) >= len(binaryCode) {
				return disassemblyInstructions, []byte{byte1, byte2}
			}
			byte3 := binaryCode[byteIndex+2]

			instructionBytes = append(instructionBytes, byte2)
			instructionBytes = append(instructionBytes, byte3)

			instructionType := proc.GetInstructionFormat34(byte2)
			if instructionType == proc.InstructionFormat3 {
				instruction.Format = proc.InstructionFormat3
			} else if instructionType == proc.InstructionFormat4 {
				// Check for incomplete instruction
				if (byteIndex + 3) >= len(binaryCode) {
					return disassemblyInstructions, []byte{byte1}
				}
				byte4 := binaryCode[byteIndex+3]

				instructionBytes = append(instructionBytes, byte4)
				instruction.Format = proc.InstructionFormat4
			}
		}
		instruction.Bytes = instructionBytes
		instruction.InstructionAddress = codeAddress.Add(relativeAddress)

		if debugGetDisassemblyInstructionsFromTextRecord {
			fmt.Printf("    Instruction bytes: % X\n", instructionBytes)
			fmt.Printf("    Instruction address: %s\n", instruction.InstructionAddress.StringHex())
		}

		byteIndex += len(instructionBytes)
		for i := 0; i < len(instructionBytes); i++ {
			relativeAddress = relativeAddress.Add(units.Int24{0x00, 0x00, 0x01})
		}

		disassemblyInstructions[instruction.InstructionAddress] = instruction
	}

	return disassemblyInstructions, []byte{}
}
