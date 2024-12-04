package core

import (
	"fmt"
	"sicsimgo/core/base"
	"sicsimgo/core/units"
)

/*
DEFINITIONS
*/
type Instruction struct {
	Format InstructionFormat
	Bytes  []byte
	Opcode Opcode
}

/*
DEBUG
*/
const debugFetchNextInstruction bool = false
const debugExecuteNextInstruction bool = false

/*
OPERATIONS
*/
func FetchNextInstruction(updatePC bool) Instruction {
	pc := base.GetRegisterPC()

	byte1 := base.GetByte(pc)
	opcode := GetOpcode(byte1)

	instructionFormat, err := GetInstructionFormat(byte1)
	if err != nil {
		// Invalid opcode
		panic(err)
	}
	if debugFetchNextInstruction {
		fmt.Printf("Fetch Next Instruction: Opcode %02X - %s\n", byte(opcode), opcode.String())
	}

	instructionBytes := []byte{byte1}

	if instructionFormat == InstructionFormat1 {
		if updatePC {
			base.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x01}))
		}
	} else if instructionFormat == InstructionFormat2 {
		byte2 := base.GetByte(pc.Add(units.Int24{0x00, 0x00, 0x01}))
		instructionBytes = append(instructionBytes, byte2)

		if updatePC {
			base.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x02}))
		}
	} else if instructionFormat == InstructionFormatSIC {
		byte2 := base.GetByte(pc.Add(units.Int24{0x00, 0x00, 0x01}))
		byte3 := base.GetByte(pc.Add(units.Int24{0x00, 0x00, 0x02}))
		instructionBytes = append(instructionBytes, byte2)
		instructionBytes = append(instructionBytes, byte3)

		if updatePC {
			base.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x03}))
		}
	} else if instructionFormat == InstructionFormat34 {
		byte2 := base.GetByte(pc.Add(units.Int24{0x00, 0x00, 0x01}))
		byte3 := base.GetByte(pc.Add(units.Int24{0x00, 0x00, 0x02}))
		instructionBytes = append(instructionBytes, byte2)
		instructionBytes = append(instructionBytes, byte3)

		instructionType := GetInstructionFormat34(byte2)
		if instructionType == InstructionFormat3 {
			instructionFormat = InstructionFormat3

			if updatePC {
				base.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x03}))
			}
		} else {
			instructionFormat = InstructionFormat4

			byte4 := base.GetByte(pc.Add(units.Int24{0x00, 0x00, 0x03}))
			instructionBytes = append(instructionBytes, byte4)

			if updatePC {
				base.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x04}))
			}
		}
	} else {
		return Instruction{
			Format: InstructionUnknown,
			Bytes:  []byte{},
			Opcode: 0x00,
		}
	}

	return Instruction{
		Format: instructionFormat,
		Bytes:  instructionBytes,
		Opcode: opcode,
	}
}

func ExecuteNextInstruction() {
	instruction := FetchNextInstruction(true)

	// halt J halt -> Stop execution
	_, address := instruction.GetOperandAddress(base.GetRegisterPC())
	pcOfInstruction := base.GetRegisterPC()
	for i := 0; i < len(instruction.Bytes); i++ {
		pcOfInstruction = pcOfInstruction.Sub(units.Int24{0x00, 0x00, 0x01})
	}
	if debugExecuteNextInstruction {
		fmt.Printf("Check for HALT: %s : %s\n", address.StringHex(), pcOfInstruction.StringHex())
	}
	if instruction.Opcode == J && address.Compare(pcOfInstruction) == 0 {
		if debugExecuteNextInstruction {
			fmt.Println("HALT")
		}
		StopSim()
	} else {
		pc := base.GetRegisterPC()
		UpdateProcState(instruction, pc)
	}

	instruction.Execute()
}
