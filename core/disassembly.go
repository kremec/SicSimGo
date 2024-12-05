package core

import (
	"fmt"
	"sicsimgo/core/units"
	"sort"
)

/*
IMPLEMENTATION
*/
var Disassembly map[units.Int24]Instruction
var InstructionList []Instruction

/*
DEBUG
*/
const debugGetDisassemblyInstructionsFromTextRecord bool = false

/*
OPERATIONS
*/
func GetInstructions(codeAddress units.Int24, code []byte) (map[units.Int24]Instruction, []byte) {
	disassemblyInstructions := make(map[units.Int24]Instruction)

	byteIndex := 0
	relativeAddress := units.Int24{}
	for byteIndex < len(code) {
		instruction := Instruction{}

		if debugGetDisassemblyInstructionsFromTextRecord {
			fmt.Printf("Looking at byte index %d\n", byteIndex)
		}

		byte1 := code[byteIndex]
		opcode := GetOpcode(byte1)
		instruction.Opcode = opcode

		instructionFormatFromOpcode, err := GetInstructionFormat(byte1)
		if err != nil {
			// Invalid opcode, perhaps data between instruction bytes
			if debugGetDisassemblyInstructionsFromTextRecord {
				fmt.Printf("    INVALID OPCODE: %s\n", opcode.String())
			}

			// Add byte as data to instructions and continue to next byte in code
			instruction := Instruction{
				InstructionAddress: codeAddress.Add(relativeAddress),
				Format:             InstructionUnknown,
				Bytes:              []byte{byte1},
				Directive:          DirectiveBYTE,
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
		if instructionFormatFromOpcode == InstructionFormat2 {
			// Check for incomplete instruction
			if (byteIndex + 1) >= len(code) {
				return disassemblyInstructions, []byte{byte1}
			}
			byte2 := code[byteIndex+1]

			r1, r2 := GetR1R2FromByte(byte2)
			instruction.R1 = r1
			instruction.R2 = r2

			instructionBytes = append(instructionBytes, byte2)
		} else if instructionFormatFromOpcode == InstructionFormatSIC {
			// Check for incomplete instruction
			if (byteIndex + 1) >= len(code) {
				return disassemblyInstructions, []byte{byte1}
			}
			byte2 := code[byteIndex+1]
			if (byteIndex + 2) >= len(code) {
				return disassemblyInstructions, []byte{byte1, byte2}
			}
			byte3 := code[byteIndex+2]

			instructionBytes = append(instructionBytes, byte2)
			instructionBytes = append(instructionBytes, byte3)
		} else if instructionFormatFromOpcode == InstructionFormat34 {
			// Check for incomplete instruction
			if (byteIndex + 1) >= len(code) {
				return disassemblyInstructions, []byte{byte1}
			}
			byte2 := code[byteIndex+1]
			if (byteIndex + 2) >= len(code) {
				return disassemblyInstructions, []byte{byte1, byte2}
			}
			byte3 := code[byteIndex+2]

			instructionBytes = append(instructionBytes, byte2)
			instructionBytes = append(instructionBytes, byte3)

			instructionType := GetInstructionFormat34(byte2)
			if instructionType == InstructionFormat3 {
				instruction.Format = InstructionFormat3
			} else if instructionType == InstructionFormat4 {
				// Check for incomplete instruction
				if (byteIndex + 3) >= len(code) {
					return disassemblyInstructions, []byte{byte1}
				}
				byte4 := code[byteIndex+3]

				instructionBytes = append(instructionBytes, byte4)
				instruction.Format = InstructionFormat4
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

		nextInstructionAddress := codeAddress.Add(relativeAddress)
		if instructionFormatFromOpcode == InstructionFormat34 {
			operand, address := instruction.GetOperandAddress(nextInstructionAddress)
			instruction.Address = address
			instruction.Operand = operand
		}

		disassemblyInstructions[instruction.InstructionAddress] = instruction
	}

	return disassemblyInstructions, []byte{}
}

func UpdateDisassemblyInstructionList() {
	adresses := units.Int24Slice{}
	for key := range Disassembly {
		adresses = append(adresses, key)
	}
	sort.Sort(adresses)

	instructionList := make([]Instruction, 0, len(Disassembly))
	for _, key := range adresses {
		instructionList = append(instructionList, Disassembly[key])
	}

	InstructionList = instructionList
}
