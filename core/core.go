package core

import (
	"fmt"
	"sicsimgo/core/base"
	"sicsimgo/core/loader"
	"sicsimgo/core/loader/bytecode"
	"sicsimgo/core/proc"
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
func GetNextDisassemblyInstruction(updatePC bool) (proc.Instruction, error) {
	pc := base.GetRegisterPC()

	if len(loader.Disassembly) == 0 {
		return proc.UnknownInstruction, loader.ErrDisassemblyEmpty()
	}
	instruction, exists := loader.Disassembly[pc]
	// Instruction not found - disassembly incorrect
	if !exists {
		if debugGetNextDisassemblyInstruction {
			fmt.Println("Instruction not found - incorrect disassembly")
		}
		// Delete all instructions after PC
		for addr := range loader.Disassembly {
			if addr.Compare(pc) > 0 {
				// Delete instructions below PC
				delete(loader.Disassembly, addr)
			}
		}

		// Replace 1st instruction above PC (which now has max address) with its bytes until PC
		maxAddr := units.Int24{0x00, 0x00, 0x00}
		for addr := range loader.Disassembly {
			if addr.Compare(maxAddr) > 0 {
				maxAddr = addr
			}
		}
		var unknownBytes []byte
		unknownBytesAddr := maxAddr
		i := 0
		for unknownBytesAddr.Compare(pc) < 0 {
			unknownBytes = append(unknownBytes, loader.Disassembly[maxAddr].Bytes[i])
			i++
			unknownBytesAddr = unknownBytesAddr.Add(units.Int24{0x00, 0x00, 0x01})
		}
		unknownBytesInstruction := proc.Instruction{
			Format:             proc.InstructionUnknown,
			Directive:          proc.DirectiveBYTE,
			Bytes:              unknownBytes,
			InstructionAddress: maxAddr,
		}
		loader.Disassembly[maxAddr] = unknownBytesInstruction

		// Disassemble from PC to last instruction byte address
		if debugGetNextDisassemblyInstruction {
			fmt.Println("Disassembling code from PC to LastInstructionByteAddress:")
		}
		codeAfterPC := base.GetSlice(pc, loader.LastInstructionByteAddress)
		instructions, bytesFromIncompleteInstruction := bytecode.GetInstructionsFromBinary(pc, codeAfterPC)
		for address, instruction := range instructions {
			if debugGetNextDisassemblyInstruction {
				fmt.Printf("    Address: %s, Format: %s, Bytes: % X, Opcode: %s, Operand: %s\n", address.StringHex(), instruction.Format.String(), instruction.Bytes, instruction.Opcode.String(), instruction.Operand.StringHex())
			}
			loader.Disassembly[address] = instruction
		}
		if len(bytesFromIncompleteInstruction) > 0 {
			// Add incomplete instruction bytes to the end of Disassembly
			addrLeftoverBytes := loader.LastInstructionByteAddress
			for i := 0; i < len(bytesFromIncompleteInstruction); i++ {
				addrLeftoverBytes.Sub(units.Int24{0x00, 0x00, 0x01})
			}
			loader.Disassembly[addrLeftoverBytes] = proc.Instruction{
				Format:             proc.InstructionUnknown,
				Bytes:              bytesFromIncompleteInstruction,
				InstructionAddress: addrLeftoverBytes,
			}
		}

		loader.UpdateInstructionList()
		UpdateProcState(base.GetRegisterPC())
		return GetNextDisassemblyInstruction(updatePC)
	}

	switch instruction.Format {
	case proc.InstructionFormat1:
		pc = pc.Add(units.Int24{0x00, 0x00, 0x01})
	case proc.InstructionFormat2:
		pc = pc.Add(units.Int24{0x00, 0x00, 0x02})
	case proc.InstructionFormatSIC:
		pc = pc.Add(units.Int24{0x00, 0x00, 0x03})
	case proc.InstructionFormat3:
		pc = pc.Add(units.Int24{0x00, 0x00, 0x03})
	case proc.InstructionFormat4:
		pc = pc.Add(units.Int24{0x00, 0x00, 0x04})
	}
	if updatePC {
		base.SetRegisterPC(pc)
	}

	// Update operand and address values
	if instruction.Format == proc.InstructionFormat3 || instruction.Format == proc.InstructionFormat4 {
		operand, address, _, _, _ := instruction.GetOperandAddress(pc)
		instruction.Operand = operand
		instruction.Address = address
	}

	return instruction, nil
}

func ExecuteNextInstruction() {
	instruction, err := GetNextDisassemblyInstruction(true)
	if err != nil {
		if err == loader.ErrDisassemblyEmpty() {
			return
		}
	}

	instruction.Execute()
	UpdateProcState(base.GetRegisterPC())

	// halt J halt -> Stop execution
	if debugExecuteNextInstruction {
		fmt.Printf("Check for HALT: %s : %s\n", instruction.InstructionAddress.StringHex(), base.GetRegisterPC().StringHex())
	}
	if instruction.Opcode == proc.J && instruction.InstructionAddress.Compare(base.GetRegisterPC()) == 0 {
		if debugExecuteNextInstruction {
			fmt.Println("HALT")
		}
		StopSim()
		return
	}
}
