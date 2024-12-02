package proc

import (
	"errors"
	"fmt"
	"sicsimgo/core"
	"sicsimgo/core/units"
)

type Instruction struct {
	Format     InstructionFormat
	Bytes      []byte
	Opcode     Opcode
	OpcodeByte byte
}

func FetchNextInstruction(updatePC bool) Instruction {
	pc := core.GetRegisterPC()

	byte1 := core.GetByte(pc)
	opcodeByte := byte1 & 0xFC
	opcode := Opcode(opcodeByte)

	instructionFormatFromOpcode, err := GetInstructionFormatFromOpcode(opcode, byte1)
	if err != nil {
		// Invalid opcode
		panic(err)
	}

	if instructionFormatFromOpcode == InstructionFormat1 {

		if updatePC {
			core.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x01}))
		}
		return Instruction{
			Format:     InstructionFormat1,
			Bytes:      []byte{byte1},
			Opcode:     opcode,
			OpcodeByte: opcodeByte,
		}
	} else if instructionFormatFromOpcode == InstructionFormat2 {
		byte2 := core.GetByte(pc.Add(units.Int24{0x00, 0x00, 0x01}))

		if updatePC {
			core.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x02}))
		}
		return Instruction{
			Format:     InstructionFormat2,
			Bytes:      []byte{byte1, byte2},
			Opcode:     opcode,
			OpcodeByte: opcodeByte,
		}
	} else if instructionFormatFromOpcode == InstructionFormatSIC {
		byte2 := core.GetByte(pc.Add(units.Int24{0x00, 0x00, 0x01}))
		byte3 := core.GetByte(pc.Add(units.Int24{0x00, 0x00, 0x02}))

		if updatePC {
			core.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x03}))
		}
		return Instruction{
			Format:     InstructionFormatSIC,
			Bytes:      []byte{byte1, byte2, byte3},
			Opcode:     opcode,
			OpcodeByte: opcodeByte,
		}
	} else if instructionFormatFromOpcode == InstructionFormat34 {
		byte2 := core.GetByte(pc.Add(units.Int24{0x00, 0x00, 0x01}))
		byte3 := core.GetByte(pc.Add(units.Int24{0x00, 0x00, 0x02}))

		e := (byte2 & 0b00010000) > 0
		if !e {
			if updatePC {
				core.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x03}))
			}
			return Instruction{
				Format:     InstructionFormat3,
				Bytes:      []byte{byte1, byte2, byte3},
				Opcode:     opcode,
				OpcodeByte: opcodeByte,
			}
		} else {
			byte4 := core.GetByte(pc.Add(units.Int24{0x00, 0x00, 0x03}))

			if updatePC {
				core.SetRegisterPC(pc.Add(units.Int24{0x00, 0x00, 0x04}))
			}
			return Instruction{
				Format:     InstructionFormat4,
				Bytes:      []byte{byte1, byte2, byte3, byte4},
				Opcode:     opcode,
				OpcodeByte: opcodeByte,
			}
		}

	} else {
		return Instruction{
			Format:     InstructionUnknown,
			Bytes:      []byte{},
			Opcode:     0x00,
			OpcodeByte: 0x00,
		}
	}
}

var debugExecuteInstruction bool = false

func (instruction Instruction) Execute() error {

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

func getNIXBPEBits(instruction []byte) (n, i, x, b, p, e bool) {
	n = instruction[0]&0b00000010 > 0
	i = instruction[0]&0b00000001 > 0
	x = instruction[1]&0b10000000 > 0
	b = instruction[1]&0b01000000 > 0
	p = instruction[1]&0b00100000 > 0
	e = instruction[1]&0b00010000 > 0

	return n, i, x, b, p, e
}

func compareOperation(r1, r2 units.Int24) {
	compareRes := r1.Compare(r2)
	if compareRes == -1 {
		core.SetRegisterSW(units.Int24{0x00, 0x00, 0x00})
	} else if compareRes == 0 {
		core.SetRegisterSW(units.Int24{0x40, 0x00, 0x00})
	} else {
		core.SetRegisterSW(units.Int24{0x80, 0x00, 0x00})
	}
}
func getCompareFromSW(sw units.Int24) int {
	return sw.Compare(units.Int24{0x00, 0x00, 0x00})
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

var debugExecuteFormat2 bool = false

func executeFormat2(instruction Instruction) error {

	r1Id := core.RegisterId(instruction.Bytes[1] & 0xF0 >> 4)
	r2Id := core.RegisterId(instruction.Bytes[1] & 0x0F)
	r1, err := core.GetRegister(r1Id)
	if err != nil {
		return err
	}
	r2, err := core.GetRegister(r2Id)
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
		core.SetRegister(r2Id, r2.Add(r1))
	case CLEAR:
		core.SetRegister(r2Id, units.Int24{0x00, 0x00, 0x00})
	case COMPR:
		compareOperation(r1, r2)
	case DIVR:
		core.SetRegister(r2Id, r2.Div(r1))
	case MULR:
		core.SetRegister(r2Id, r2.Mul(r1))
	case RMO:
		core.SetRegister(r2Id, r1)
	// TODO: SHIFT OPERATIONS
	case SHIFTL:
	case SHIFTR:
	case SUBR:
		core.SetRegister(r2Id, r2.Sub(r1))
	// TODO: SYSCALL
	case SVC:
	case TIXR:
		core.SetRegisterX(core.GetRegisterX().Add(units.Int24{0x00, 0x00, 0x01}))
		compareOperation(core.GetRegisterX(), r1)
	}

	return nil
}

var debugGetOperandAddress bool = false

func (instruction Instruction) getOperandAddress(pc units.Int24) (units.Int24, units.Int24) {
	n, i, x, b, p, _ := getNIXBPEBits(instruction.Bytes)
	relativeAddressingMode, err := GetRelativeAdressingModes(b, p)
	if err != nil {
		// Invalid addressing
		return units.Int24{}, units.Int24{}
	}
	absoluteAddressingMode := GetAbsoluteAdressingModes(n, i)
	indexAddressingMode := GetIndexAdressingModes(x)

	var address units.Int24
	switch instruction.Format {
	case InstructionFormatSIC:
		address = units.Int24{0x00, instruction.Bytes[1] & 0b01111111, instruction.Bytes[2]}
	case InstructionFormat3:
		// Sign-extend operand
		if (instruction.Bytes[1] & 0b00001000) > 0 {
			address = units.Int24{0xFF, (instruction.Bytes[1] & 0b00001111) | 0b11110000, instruction.Bytes[2]}
		} else {
			address = units.Int24{0x00, instruction.Bytes[1] & 0b00001111, instruction.Bytes[2]}
		}
	case InstructionFormat4:
		address = units.Int24{instruction.Bytes[1] & 0b00001111, instruction.Bytes[2], instruction.Bytes[3]}
	}

	if debugGetOperandAddress {
		fmt.Println("    Instruction:", instruction.Opcode.String())
		fmt.Println("    Address: ", address.StringHex())
		fmt.Println("        Relative addressing mode: ", relativeAddressingMode.String())
		fmt.Println("        Absolute addressing mode: ", absoluteAddressingMode.String())
	}

	// Relative addressing
	switch relativeAddressingMode {
	case DirectRelativeAddressing:
		// Do nothing
	case PCRelativeAddressing:
		address = address.Add(pc)
	case BaseRelativeAddressing:
		address = address.Add(core.GetRegisterB()) // TODO: Check if ok!
	}

	if debugGetOperandAddress {
		fmt.Println("    Address after relative addressing: ", address.StringHex())
		fmt.Println("        PC: ", pc.StringHex())
	}

	// Index addressing
	if indexAddressingMode {
		address = address.Add(core.GetRegisterX())
	}

	if debugGetOperandAddress {
		fmt.Println("    Address after index addressing: ", address.StringHex())
	}

	var operand units.Int24
	// Absolute addressing
	switch absoluteAddressingMode {
	case SICAbsoluteAddressing:
		operand = core.GetWord(address)
	case ImmediateAbsoluteAddressing:
		operand = address
	case IndirectAbsoluteAddressing:
		operand = core.GetWord(core.GetWord(address))
	case DirectAbsoluteAddressing:
		operand = core.GetWord(address)
	}

	if debugGetOperandAddress {
		fmt.Println("    Operand:", operand.StringHex())
		fmt.Println()
	}

	return operand, address
}
func executeFormatSIC34(instruction Instruction) error {
	operand, address := instruction.getOperandAddress(core.GetRegisterPC())

	switch instruction.Opcode {
	case ADD:
		core.SetRegisterA(core.GetRegisterA().Add(operand))
	// TODO: FLOAT
	case ADDF:
	case AND:
		core.SetRegisterA(core.GetRegisterA().And(operand))
	case COMP:
		compareOperation(core.GetRegisterA(), operand)
	// TODO: FLOAT
	case COMPF:
	case DIV:
		core.SetRegisterA(core.GetRegisterA().Div(operand))
	// TODO: FLOAT
	case DIVF:
	case J:
		core.SetRegisterPC(address)
	case JEQ:
		if getCompareFromSW(core.GetRegisterSW()) == 0 {
			core.SetRegisterPC(address)
		}
	case JGT:
		if getCompareFromSW(core.GetRegisterSW()) == 1 {
			core.SetRegisterPC(address)
		}
	case JLT:
		if getCompareFromSW(core.GetRegisterSW()) == -1 {
			core.SetRegisterPC(address)
		}
	case JSUB:
		core.SetRegisterL(core.GetRegisterPC())
		core.SetRegisterPC(address)
	case LDA:
		core.SetRegisterA(operand)
	case LDB:
		core.SetRegisterB(operand)
	case LDCH:
		core.SetRegisterA(units.Int24{0x00, 0x00, operand[2]})
	// TODO: FLOAT
	case LDF:
	case LDL:
		core.SetRegisterL(operand)
	case LDS:
		core.SetRegisterS(operand)
	case LDT:
		core.SetRegisterT(operand)
	case LDX:
		core.SetRegisterX(operand)
	// TODO: SYSCALL
	case LPS:
	case MUL:
		core.SetRegisterA(core.GetRegisterA().Mul(operand))
	// TODO: FLOAT
	case MULF:
	case OR:
		core.SetRegisterA(core.GetRegisterA().Or(operand))
	case RD:
		if readByte, err := core.Read(core.Device(operand[0])); err == nil {
			core.SetRegisterA(units.Int24{0x00, 0x00, readByte})
		}
	case RSUB:
		core.SetRegisterPC(core.GetRegisterL())
	// TODO: SYSCALL
	case SSK:
	case STA:
		core.SetWord(operand, core.GetRegisterA())
	case STB:
		core.SetWord(operand, core.GetRegisterB())
	case STCH:
		core.SetByte(operand, core.GetRegisterA()[2])
	// TODO: FLOAT
	case STF:
	// TODO: SYSCALL
	case STI:
	case STL:
		core.SetWord(operand, core.GetRegisterL())
	case STS:
		core.SetWord(operand, core.GetRegisterS())
	case STSW:
		core.SetWord(operand, core.GetRegisterSW())
	case STT:
		core.SetWord(operand, core.GetRegisterT())
	case STX:
		core.SetWord(operand, core.GetRegisterX())
	case SUB:
		core.SetRegisterA(core.GetRegisterA().Sub(operand))
	// TODO: FLOAT
	case SUBF:
	// TODO: SYSCALL
	case TD:
	case TIX:
		core.SetRegisterX(core.GetRegisterX().Add(units.Int24{0x00, 0x00, 0x01}))
		compareOperation(core.GetRegisterX(), operand)
	case WD:
		if err := core.Write(core.Device(operand[2]), core.GetRegisterA()[2]); err != nil {
			// TODO: Handle error
		}
	}

	return nil
}
