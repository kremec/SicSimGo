package core

import (
	"errors"
	"fmt"
	"sicsimgo/core/base"
	"sicsimgo/core/units"
)

/*
DEBUG
*/
var debugExecuteInstruction bool = false
var debugExecuteFormat2 bool = false

/*
OPERATIONS
*/
func compareOperation(r1, r2 units.Int24) {
	compareRes := r1.Compare(r2)
	if compareRes == -1 {
		base.SetRegisterSW(units.Int24{0x00, 0x00, 0x00})
	} else if compareRes == 0 {
		base.SetRegisterSW(units.Int24{0x40, 0x00, 0x00})
	} else {
		base.SetRegisterSW(units.Int24{0x80, 0x00, 0x00})
	}
}
func getCompareFromSW(sw units.Int24) int {
	return sw.Compare(units.Int24{0x00, 0x00, 0x00})
}

func (instruction Instruction) Execute() error {

	if debugExecuteInstruction {
		fmt.Printf("Execute Instruction: Opcode %02X - Format %d\n", instruction.Opcode, instruction.Format)
	}

	switch instruction.Format {
	case InstructionFormat1:
		if debugExecuteInstruction {
			fmt.Printf("Instruction: Opcode %02X - Format %d - Bytes [%02X]\n", instruction.Opcode, instruction.Format, instruction.Bytes[0])
		}
		return executeFormat1(instruction)
	case InstructionFormat2:
		if debugExecuteInstruction {
			fmt.Printf("Instruction: Opcode %02X - Format %d - Bytes [%02X %02X]\n", instruction.Opcode, instruction.Format, instruction.Bytes[0], instruction.Bytes[1])
		}
		return executeFormat2(instruction)
	case InstructionFormatSIC:
		if debugExecuteInstruction {
			fmt.Printf("Instruction: Opcode %02X - Format %d - Bytes [%02X %02X %02X]\n", instruction.Opcode, instruction.Format, instruction.Bytes[0], instruction.Bytes[1], instruction.Bytes[2])
		}
		return executeFormatSIC34(instruction)
	case InstructionFormat3:
		if debugExecuteInstruction {
			fmt.Printf("Instruction: Opcode %02X - Format %d - Bytes [%02X %02X %02X]\n", instruction.Opcode, instruction.Format, instruction.Bytes[0], instruction.Bytes[1], instruction.Bytes[2])
		}
		return executeFormatSIC34(instruction)
	case InstructionFormat4:
		if debugExecuteInstruction {
			fmt.Printf("Instruction: Opcode %02X - Format %d - Bytes [%02X %02X %02X %02X]\n", instruction.Opcode, instruction.Format, instruction.Bytes[0], instruction.Bytes[1], instruction.Bytes[2], instruction.Bytes[3])
		}
		return executeFormatSIC34(instruction)
	default:
		return errors.New("Invalid instruction format")
	}
}

func executeFormat1(instruction Instruction) error {
	switch instruction.Opcode {
	// TODO: FLOAT
	case FIX:
	case FLOAT:
	// TODO: SYSCALL
	case HIO:
	case NORM:
	case SIO:
	case TIO:
	}

	return nil
}

func executeFormat2(instruction Instruction) error {

	r1Id, r2Id := GetR1R2FromByte(instruction.Bytes[1])
	r1, err := r1Id.GetRegister()
	if err != nil {
		return err
	}
	r2, err := r2Id.GetRegister()
	if err != nil {
		return err
	}

	if debugExecuteFormat2 {
		fmt.Printf("    Instruction bytes: %08b %08b\n", instruction.Bytes[0], instruction.Bytes[1])
		fmt.Println("    Register 1: ", r1Id.String())
		fmt.Println("        Value: ", r1.StringHex())
		fmt.Println("    Register 2: ", r2Id.String())
		fmt.Println("        Value: ", r2.StringHex())
	}

	switch instruction.Opcode {
	case ADDR:
		r2Id.SetRegister(r2.Add(r1))
	case CLEAR:
		r2Id.SetRegister(units.Int24{0x00, 0x00, 0x00})
	case COMPR:
		compareOperation(r1, r2)
	case DIVR:
		r2Id.SetRegister(r2.Div(r1))
	case MULR:
		r2Id.SetRegister(r2.Mul(r1))
	case RMO:
		r2Id.SetRegister(r1)
	// TODO: SHIFT OPERATIONS
	case SHIFTL:
	case SHIFTR:
	case SUBR:
		r2Id.SetRegister(r2.Sub(r1))
	// TODO: SYSCALL
	case SVC:
	case TIXR:
		base.SetRegisterX(base.GetRegisterX().Add(units.Int24{0x00, 0x00, 0x01}))
		compareOperation(base.GetRegisterX(), r1)
	}

	return nil
}

func executeFormatSIC34(instruction Instruction) error {
	operand, address := instruction.GetOperandAddress(base.GetRegisterPC())

	switch instruction.Opcode {
	case ADD:
		base.SetRegisterA(base.GetRegisterA().Add(operand))
	// TODO: FLOAT
	case ADDF:
	case AND:
		base.SetRegisterA(base.GetRegisterA().And(operand))
	case COMP:
		compareOperation(base.GetRegisterA(), operand)
	// TODO: FLOAT
	case COMPF:
	case DIV:
		base.SetRegisterA(base.GetRegisterA().Div(operand))
	// TODO: FLOAT
	case DIVF:
	case J:
		base.SetRegisterPC(address)
	case JEQ:
		if getCompareFromSW(base.GetRegisterSW()) == 0 {
			base.SetRegisterPC(address)
		}
	case JGT:
		if getCompareFromSW(base.GetRegisterSW()) == 1 {
			base.SetRegisterPC(address)
		}
	case JLT:
		if getCompareFromSW(base.GetRegisterSW()) == -1 {
			base.SetRegisterPC(address)
		}
	case JSUB:
		base.SetRegisterL(base.GetRegisterPC())
		base.SetRegisterPC(address)
	case LDA:
		base.SetRegisterA(operand)
	case LDB:
		base.SetRegisterB(operand)
	case LDCH:
		base.SetRegisterA(units.Int24{0x00, 0x00, operand[2]})
	// TODO: FLOAT
	case LDF:
	case LDL:
		base.SetRegisterL(operand)
	case LDS:
		base.SetRegisterS(operand)
	case LDT:
		base.SetRegisterT(operand)
	case LDX:
		base.SetRegisterX(operand)
	// TODO: SYSCALL
	case LPS:
	case MUL:
		base.SetRegisterA(base.GetRegisterA().Mul(operand))
	// TODO: FLOAT
	case MULF:
	case OR:
		base.SetRegisterA(base.GetRegisterA().Or(operand))
	case RD:
		if readByte, err := base.Read(base.Device(operand[0])); err == nil {
			base.SetRegisterA(units.Int24{0x00, 0x00, readByte})
		}
	case RSUB:
		base.SetRegisterPC(base.GetRegisterL())
	// TODO: SYSCALL
	case SSK:
	case STA:
		base.SetWord(operand, base.GetRegisterA())
	case STB:
		base.SetWord(operand, base.GetRegisterB())
	case STCH:
		base.SetByte(operand, base.GetRegisterA()[2])
	// TODO: FLOAT
	case STF:
	// TODO: SYSCALL
	case STI:
	case STL:
		base.SetWord(operand, base.GetRegisterL())
	case STS:
		base.SetWord(operand, base.GetRegisterS())
	case STSW:
		base.SetWord(operand, base.GetRegisterSW())
	case STT:
		base.SetWord(operand, base.GetRegisterT())
	case STX:
		base.SetWord(operand, base.GetRegisterX())
	case SUB:
		base.SetRegisterA(base.GetRegisterA().Sub(operand))
	// TODO: FLOAT
	case SUBF:
	// TODO: SYSCALL
	case TD:
	case TIX:
		base.SetRegisterX(base.GetRegisterX().Add(units.Int24{0x00, 0x00, 0x01}))
		compareOperation(base.GetRegisterX(), operand)
	case WD:
		if err := base.Write(base.Device(operand[2]), base.GetRegisterA()[2]); err != nil {
			// TODO: Handle error
		}
	}

	return nil
}
