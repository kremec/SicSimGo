package core

import (
	"fmt"
	"sicsimgo/core/base"
	"sicsimgo/core/units"
)

/*
DEBUG
*/
const debugGetNextDisassemblyInstruction bool = false
const debugExecuteNextInstruction bool = false

/*
OPERATIONS
*/
func GetNextDisassemblyInstruction(updatePC bool) (Instruction, error) {
	pc := base.GetRegisterPC()

	if len(Disassembly) == 0 {
		return UnknownInstruction, ErrDisassemblyEmpty()
	}
	instruction, exists := Disassembly[pc]
	// Instruction not found - disassembly incorrect
	if !exists {
		if debugGetNextDisassemblyInstruction {
			fmt.Println("Instruction not found - incorrect disassembly")
		}
		// Delete all instructions after PC
		for addr := range Disassembly {
			if addr.Compare(pc) > 0 {
				// Delete instructions below PC
				delete(Disassembly, addr)
			}
		}

		// Replace 1st instruction above PC (which now has max address) with its bytes until PC
		maxAddr := units.Int24{0x00, 0x00, 0x00}
		for addr := range Disassembly {
			if addr.Compare(maxAddr) > 0 {
				maxAddr = addr
			}
		}
		var unknownBytes []byte
		unknownBytesAddr := maxAddr
		i := 0
		for unknownBytesAddr.Compare(pc) < 0 {
			unknownBytes = append(unknownBytes, Disassembly[maxAddr].Bytes[i])
			i++
			unknownBytesAddr = unknownBytesAddr.Add(units.Int24{0x00, 0x00, 0x01})
		}
		unknownBytesInstruction := Instruction{
			Format:             InstructionUnknown,
			Directive:          DirectiveBYTE,
			Bytes:              unknownBytes,
			InstructionAddress: maxAddr,
		}
		Disassembly[maxAddr] = unknownBytesInstruction

		// Disassemble from PC to last instruction byte address
		if debugGetNextDisassemblyInstruction {
			fmt.Println("Disassembling code from PC to LastInstructionByteAddress:")
		}
		codeAfterPC := base.GetSlice(pc, LastInstructionByteAddress)
		instructions, bytesFromIncompleteInstruction := GetInstructions(pc, codeAfterPC)
		for address, instruction := range instructions {
			if debugLoadProgram {
				fmt.Printf("    Address: %s, Format: %s, Bytes: % X, Opcode: %s, Operand: %s\n", address.StringHex(), instruction.Format.String(), instruction.Bytes, instruction.Opcode.String(), instruction.Operand.StringHex())
			}
			Disassembly[address] = instruction
		}
		if len(bytesFromIncompleteInstruction) > 0 {
			// Add incomplete instruction bytes to the end of Disassembly
			addrLeftoverBytes := LastInstructionByteAddress
			for i := 0; i < len(bytesFromIncompleteInstruction); i++ {
				addrLeftoverBytes.Sub(units.Int24{0x00, 0x00, 0x01})
			}
			Disassembly[addrLeftoverBytes] = Instruction{
				Format:             InstructionUnknown,
				Bytes:              bytesFromIncompleteInstruction,
				InstructionAddress: addrLeftoverBytes,
			}
		}

		UpdateDisassemblyInstructionList()
		UpdateProcState(base.GetRegisterPC())
		return GetNextDisassemblyInstruction(updatePC)
	}
	if updatePC {
		switch instruction.Format {
		case InstructionFormat1:
			base.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x01}))
		case InstructionFormat2:
			base.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x02}))
		case InstructionFormatSIC:
			base.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x03}))
		case InstructionFormat3:
			base.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x03}))
		case InstructionFormat4:
			base.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x04}))
		}
	}

	// Update operand and address values
	if instruction.Format == InstructionFormat3 || instruction.Format == InstructionFormat4 {
		operand, address := instruction.GetOperandAddress(base.GetRegisterPC())
		instruction.Address = address
		instruction.Operand = operand
	}

	return instruction, nil
}

func ExecuteNextInstruction() {
	instruction, err := GetNextDisassemblyInstruction(true)
	if err != nil {
		if err == ErrDisassemblyEmpty() {
			return
		}
	}

	instruction.Execute()
	UpdateProcState(base.GetRegisterPC())

	// halt J halt -> Stop execution
	if debugExecuteNextInstruction {
		fmt.Printf("Check for HALT: %s : %s\n", instruction.InstructionAddress.StringHex(), base.GetRegisterPC().StringHex())
	}
	if instruction.Opcode == J && instruction.InstructionAddress.Compare(base.GetRegisterPC()) == 0 {
		if debugExecuteNextInstruction {
			fmt.Println("HALT")
		}
		StopSim()
		return
	}
}
