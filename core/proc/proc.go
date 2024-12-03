package proc

import (
	"fmt"
	"sicsimgo/core"
	"sicsimgo/core/units"
)

type ProcState struct {
	Instruction      Instruction
	N, I, X, B, P, E bool
	Operand          units.Int24
	Address          units.Int24
}

func InitProcState() {
	instruction := FetchNextInstruction(false)
	UpdateProcState(instruction, core.GetRegisterPC())
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
	var n, i, x, b, p, e = getNIXBPEBits(instruction.Bytes)

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

var CurrentProcState ProcState = ProcState{}

func ExecuteNextInstruction() {
	instruction := FetchNextInstruction(true)

	// halt J halt -> Stop execution
	endOfProgram := false
	_, address := instruction.GetOperandAddress(core.GetRegisterPC())
	pcOfInstruction := core.GetRegisterPC()
	for i := 0; i < len(instruction.Bytes); i++ {
		pcOfInstruction = pcOfInstruction.Sub(units.Int24{0x00, 0x00, 0x01})
	}
	fmt.Printf("Check for HALT: %s : %s\n", address.StringHex(), pcOfInstruction.StringHex())
	if instruction.Opcode == J && address.Compare(pcOfInstruction) == 0 {
		// End of program
		core.SimExecuteState = false
		endOfProgram = true
	}

	if !endOfProgram {
		pc := core.GetRegisterPC()
		UpdateProcState(instruction, pc)
	}
	instruction.Execute()
}
