package core

import (
	"fmt"
	"sicsimgo/core/base"
	"sicsimgo/core/units"
)

/*
DEFINITIONS
*/
type ExecuteState bool

type ProcState struct {
	Instruction      Instruction
	N, I, X, B, P, E bool
}

const (
	ExecuteStartState ExecuteState = true
	ExecuteStopState  ExecuteState = false
)

/*
IMPLEMENTATION
*/
var SimExecuteState ExecuteState = ExecuteStopState
var CurrentProcState ProcState = ProcState{}

/*
DEBUG
*/
var debugUpdateProcState bool = true

/*
OPERATIONS
*/
func UpdateProcState(pc units.Int24) {
	var instruction Instruction

	if debugUpdateProcState {
		fmt.Println("PC: ", pc.StringHex())
	}

	// Get next instruction and PC (simulate fetch)
	nextInstruction, err := GetNextDisassemblyInstruction(false)
	if err != nil {
		if err == ErrDisassemblyEmpty() {
			return
		}
	}
	instruction = nextInstruction

	if debugUpdateProcState {
		fmt.Printf("Opcode: %s, Instruction bytes: % X\n", instruction.Opcode.String(), instruction.Bytes)
		fmt.Println()
	}

	n, i, x, b, p, e := instruction.GetNIXBPEBits()

	CurrentProcState = ProcState{
		Instruction: instruction,
		N:           n,
		I:           i,
		X:           x,
		B:           b,
		P:           p,
		E:           e,
	}
}

func StopSim() {
	SimExecuteState = ExecuteStopState
}

func ResetSim() {
	SimExecuteState = ExecuteStopState
	Disassembly = make(map[units.Int24]Instruction)
	base.ResetRegisters()
	base.ResetMemory()
}
