package core

import (
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
	Operand          units.Int24
	Address          units.Int24
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
OPERATIONS
*/
func InitProcState() {
	instruction := FetchNextInstruction(false)
	UpdateProcState(instruction, base.GetRegisterPC())
}

func UpdateProcState(instruction Instruction, pc units.Int24) {

	//fmt.Println("PC: ", pc.StringHex())

	// Get next instruction and PC (simulate fetch)
	if (pc.Compare(units.Int24{}) != 0) {
		instruction = FetchNextInstruction(false)
	}
	for i := 0; i < len(instruction.Bytes); i++ {
		pc = pc.Add(units.Int24{0x00, 0x00, 0x01})
	}

	//fmt.Println("Instruction: ", instruction.Bytes, " ", instruction.Opcode.String())
	//fmt.Println("PC after instruction: ", pc.StringHex())
	//fmt.Println("")

	var operand, address units.Int24
	if instruction.Format == InstructionFormat3 || instruction.Format == InstructionFormat4 {
		operand, address = instruction.GetOperandAddress(pc)
	}
	var n, i, x, b, p, e = GetNIXBPEBits(instruction.Bytes)

	CurrentProcState = ProcState{
		Instruction: instruction,
		N:           n,
		I:           i,
		X:           x,
		B:           b,
		P:           p,
		E:           e,
		Operand:     operand,
		Address:     address,
	}
}

func StopSim() {
	SimExecuteState = ExecuteStopState
}

func ResetSim() {
	SimExecuteState = ExecuteStopState
	Disassembly = []DisassemblyInstruction{}
	base.ResetRegisters()
	base.ResetMemory()
}
