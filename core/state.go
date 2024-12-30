package core

import (
	"fmt"
	"sicsimgo/core/base"
	"sicsimgo/core/loader"
	"sicsimgo/core/proc"
	"sicsimgo/core/units"
)

/*
DEFINITIONS
*/
type ExecuteState bool

type ProcState struct {
	Instruction      proc.Instruction
	N, I, X, B, P, E bool
}

const (
	ExecuteStartState ExecuteState = true
	ExecuteStopState  ExecuteState = false
)

/*
IMPLEMENTATION
*/
var LoadedProgramTypeState loader.LoadedProgramType = loader.None
var SimExecuteState ExecuteState = ExecuteStopState
var CurrentProcState ProcState = ProcState{}

/*
DEBUG
*/
var debugUpdateProcState bool = false

/*
OPERATIONS
*/
func UpdateProcState(pc units.Int24) {
	var instruction proc.Instruction

	if debugUpdateProcState {
		fmt.Println("PC: ", pc.StringHex())
	}

	// Get next instruction and PC (simulate fetch)
	nextInstruction, err := GetNextDisassemblyInstruction(false)
	if err != nil {
		if err == loader.ErrDisassemblyIncorrect() {
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
	CurrentProcState = ProcState{}
	loader.ResetDissasembly()
	base.ResetRegisters()
	base.ResetMemory()
}
