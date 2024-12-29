package proc

import (
	"fmt"
	"sicsimgo/core/base"
	"sicsimgo/core/units"
)

/*
DEFINITIONS
*/
type InstructionFormat int
type Opcode byte
type Directive string

type Instruction struct {
	InstructionAddress units.Int24

	Format    InstructionFormat
	Bytes     []byte
	Opcode    Opcode
	Directive Directive

	Address units.Int24
	Operand units.Int24
	R1, R2  base.RegisterId

	RelativeAddressingMode RelativeAddressingMode
	IndexAddressingMode    IndexAddressingMode
	AbsoluteAddressingMode AbsoluteAddressingMode
}

const (
	InstructionFormat1   InstructionFormat = 1
	InstructionFormat2   InstructionFormat = 2
	InstructionFormatSIC InstructionFormat = 0
	InstructionFormat3   InstructionFormat = 3
	InstructionFormat4   InstructionFormat = 4

	InstructionFormat34 InstructionFormat = 34
	InstructionUnknown  InstructionFormat = -1
)

const (
	ADD    Opcode = 0x18
	ADDF   Opcode = 0x58
	ADDR   Opcode = 0x90
	AND    Opcode = 0x40
	CLEAR  Opcode = 0xB4
	COMP   Opcode = 0x28
	COMPF  Opcode = 0x88
	COMPR  Opcode = 0xA0
	DIV    Opcode = 0x24
	DIVF   Opcode = 0x64
	DIVR   Opcode = 0x9C
	FIX    Opcode = 0xC4
	FLOAT  Opcode = 0xC0
	HIO    Opcode = 0xF4
	J      Opcode = 0x3C
	JEQ    Opcode = 0x30
	JGT    Opcode = 0x34
	JLT    Opcode = 0x38
	JSUB   Opcode = 0x48
	LDA    Opcode = 0x00
	LDB    Opcode = 0x68
	LDCH   Opcode = 0x50
	LDF    Opcode = 0x70
	LDL    Opcode = 0x08
	LDS    Opcode = 0x6C
	LDT    Opcode = 0x74
	LDX    Opcode = 0x04
	LPS    Opcode = 0xD0
	MUL    Opcode = 0x20
	MULF   Opcode = 0x60
	MULR   Opcode = 0x98
	NORM   Opcode = 0xC8
	OR     Opcode = 0x44
	RD     Opcode = 0xD8
	RMO    Opcode = 0xAC
	RSUB   Opcode = 0x4C
	SHIFTL Opcode = 0xA4
	SHIFTR Opcode = 0xA8
	SIO    Opcode = 0xF0
	SSK    Opcode = 0xEC
	STA    Opcode = 0x0C
	STB    Opcode = 0x78
	STCH   Opcode = 0x54
	STF    Opcode = 0x80
	STI    Opcode = 0xD4
	STL    Opcode = 0x14
	STS    Opcode = 0x7C
	STSW   Opcode = 0xE8
	STT    Opcode = 0x84
	STX    Opcode = 0x10
	SUB    Opcode = 0x1C
	SUBF   Opcode = 0x5C
	SUBR   Opcode = 0x94
	SVC    Opcode = 0xB0
	TD     Opcode = 0xE0
	TIO    Opcode = 0xF8
	TIX    Opcode = 0x2C
	TIXR   Opcode = 0xB8
	WD     Opcode = 0xDC
)

var opcodesFormat1 = map[Opcode]InstructionFormat{
	FIX:   InstructionFormat1,
	FLOAT: InstructionFormat1,
	HIO:   InstructionFormat1,
	NORM:  InstructionFormat1,
	SIO:   InstructionFormat1,
	TIO:   InstructionFormat1,
}
var opcodesFormat2 = map[Opcode]InstructionFormat{
	ADDR:   InstructionFormat2,
	CLEAR:  InstructionFormat2,
	COMPR:  InstructionFormat2,
	DIVR:   InstructionFormat2,
	MULR:   InstructionFormat2,
	RMO:    InstructionFormat2,
	SHIFTL: InstructionFormat2,
	SHIFTR: InstructionFormat2,
	SUBR:   InstructionFormat2,
	SVC:    InstructionFormat2,
	TIXR:   InstructionFormat2,
}
var opcodesFormat34 = map[Opcode]InstructionFormat{
	ADD:   InstructionFormat34,
	ADDF:  InstructionFormat34,
	AND:   InstructionFormat34,
	COMP:  InstructionFormat34,
	COMPF: InstructionFormat34,
	DIV:   InstructionFormat34,
	DIVF:  InstructionFormat34,
	J:     InstructionFormat34,
	JEQ:   InstructionFormat34,
	JGT:   InstructionFormat34,
	JLT:   InstructionFormat34,
	JSUB:  InstructionFormat34,
	LDA:   InstructionFormat34,
	LDB:   InstructionFormat34,
	LDCH:  InstructionFormat34,
	LDF:   InstructionFormat34,
	LDL:   InstructionFormat34,
	LDS:   InstructionFormat34,
	LDT:   InstructionFormat34,
	LDX:   InstructionFormat34,
	LPS:   InstructionFormat34,
	MUL:   InstructionFormat34,
	MULF:  InstructionFormat34,
	OR:    InstructionFormat34,
	RD:    InstructionFormat34,
	RSUB:  InstructionFormat34,
	SSK:   InstructionFormat34,
	STA:   InstructionFormat34,
	STB:   InstructionFormat34,
	STCH:  InstructionFormat34,
	STF:   InstructionFormat34,
	STI:   InstructionFormat34,
	STL:   InstructionFormat34,
	STS:   InstructionFormat34,
	STSW:  InstructionFormat34,
	STT:   InstructionFormat34,
	STX:   InstructionFormat34,
	SUB:   InstructionFormat34,
	SUBF:  InstructionFormat34,
	TD:    InstructionFormat34,
	TIX:   InstructionFormat34,
	WD:    InstructionFormat34,
}

const (
	DirectiveNone Directive = ""
	DirectiveBYTE Directive = "BYTE"
)

var UnknownInstruction Instruction = Instruction{
	Format: InstructionUnknown,
}

/*
OPERATIONS
*/
func GetOpcode(byte1 byte) Opcode {
	return Opcode(byte1 & 0xFC)
}

func GetInstructionFormat(byte1 byte) (InstructionFormat, error) {

	opcode := GetOpcode(byte1)

	// Check opcode format
	if _, exists := opcodesFormat1[opcode]; exists {
		return InstructionFormat1, nil
	}

	// Format 2
	if _, exists := opcodesFormat2[opcode]; exists {
		return InstructionFormat2, nil
	}

	// SIC format
	if (byte1 & 0b00000011) == 0x00 {
		return InstructionFormatSIC, nil
	}

	// Format 3/4
	if _, exists := opcodesFormat34[opcode]; exists {
		return InstructionFormat34, nil
	}

	return InstructionUnknown, ErrInvalidOpcode(opcode)
}

func GetInstructionFormat34(byte2 byte) InstructionFormat {
	e := (byte2 & 0b00010000) > 0
	if e {
		return InstructionFormat4
	} else {
		return InstructionFormat3
	}
}

func (instruction Instruction) IsJumpInstruction() bool {
	switch instruction.Opcode {
	case J, JSUB, JEQ, JGT, JLT:
		return true
	}
	return false
}

func (instruction Instruction) IsStoreInstruction() bool {
	switch instruction.Opcode {
	case STCH, STA, STB, STF, STSW, STT, STX:
		return true
	}
	return false
}

func (instruction Instruction) IsFormatSIC34() bool {
	return instruction.Format == InstructionFormatSIC || instruction.Format == InstructionFormat3 || instruction.Format == InstructionFormat4
}

func (instruction Instruction) GetInstructionBytes() []byte {

	toInt := func(b bool) byte {
		if b {
			return 1
		} else {
			return 0
		}
	}

	switch instruction.Format {
	case InstructionFormat1:
		byte1 := byte(instruction.Opcode)
		return []byte{byte1}
	case InstructionFormat2:
		byte1 := byte(instruction.Opcode)
		byte2 := byte((instruction.R1&0x0F)<<4 | (instruction.R2 & 0x0F))
		return []byte{byte1, byte2}
	case InstructionFormat3:
		n, i, x, b, p, e := instruction.GenerateNIXBPEBits()
		byte1 := byte(instruction.Opcode) | (byte(toInt(n)) << 1) | byte(toInt(i))
		byte2 := byte((toInt(x) << 7) | (toInt(b) << 6) | (toInt(p) << 5) | (toInt(e) << 4) | (instruction.Address[1] & 0x0F))
		byte3 := byte(instruction.Address[2])
		return []byte{byte1, byte2, byte3}
	case InstructionFormat4:
		n, i, x, b, p, e := instruction.GenerateNIXBPEBits()
		byte1 := byte(instruction.Opcode) | (byte(toInt(n)) << 1) | byte(toInt(i))
		byte2 := byte((toInt(x) << 7) | (toInt(b) << 6) | (toInt(p) << 5) | (toInt(e) << 4) | (instruction.Address[0] & 0x0F))
		byte3 := byte(instruction.Address[1])
		byte4 := byte(instruction.Address[2])
		return []byte{byte1, byte2, byte3, byte4}
	}

	return []byte{}
}

/*
STRINGS
*/
func (instructionFormat InstructionFormat) String() string {
	switch instructionFormat {
	case InstructionFormat1:
		return "Format 1"
	case InstructionFormat2:
		return "Format 2"
	case InstructionFormatSIC:
		return "SIC format"
	case InstructionFormat3:
		return "Format 3"
	case InstructionFormat4:
		return "Format 4"
	case InstructionFormat34:
		return "Format 3/4"
	case InstructionUnknown:
		return "Unknown"
	}
	return "Not implemented"
}

func (opcode Opcode) String() string {
	switch opcode {
	case ADD:
		return "ADD"
	case ADDF:
		return "ADDF"
	case ADDR:
		return "ADDR"
	case AND:
		return "AND"
	case CLEAR:
		return "CLEAR"
	case COMP:
		return "COMP"
	case COMPF:
		return "COMPF"
	case COMPR:
		return "COMPR"
	case DIV:
		return "DIV"
	case DIVF:
		return "DIVF"
	case DIVR:
		return "DIVR"
	case FIX:
		return "FIX"
	case FLOAT:
		return "FLOAT"
	case HIO:
		return "HIO"
	case J:
		return "J"
	case JEQ:
		return "JEQ"
	case JGT:
		return "JGT"
	case JLT:
		return "JLT"
	case JSUB:
		return "JSUB"
	case LDA:
		return "LDA"
	case LDB:
		return "LDB"
	case LDCH:
		return "LDCH"
	case LDF:
		return "LDF"
	case LDL:
		return "LDL"
	case LDS:
		return "LDS"
	case LDT:
		return "LDT"
	case LDX:
		return "LDX"
	case LPS:
		return "LPS"
	case MUL:
		return "MUL"
	case MULF:
		return "MULF"
	case MULR:
		return "MULR"
	case NORM:
		return "NORM"
	case OR:
		return "OR"
	case RD:
		return "RD"
	case RMO:
		return "RMO"
	case RSUB:
		return "RSUB"
	case SHIFTL:
		return "SHIFTL"
	case SHIFTR:
		return "SHIFTR"
	case SIO:
		return "SIO"
	case SSK:
		return "SSK"
	case STA:
		return "STA"
	case STB:
		return "STB"
	case STCH:
		return "STCH"
	case STF:
		return "STF"
	case STI:
		return "STI"
	case STL:
		return "STL"
	case STS:
		return "STS"
	case STSW:
		return "STSW"
	case STT:
		return "STT"
	case STX:
		return "STX"
	case SUB:
		return "SUB"
	case SUBF:
		return "SUBF"
	case SUBR:
		return "SUBR"
	case SVC:
		return "SVC"
	case TD:
		return "TD"
	case TIO:
		return "TIO"
	case TIX:
		return "TIX"
	case TIXR:
		return "TIXR"
	case WD:
		return "WD"
	}
	return "Unknown"
}

func (instruction Instruction) String() string {
	if instruction.Directive == DirectiveBYTE {
		return fmt.Sprintf("%s : %s\n    Bytes: %X",
			instruction.InstructionAddress.StringHex(), instruction.Directive,
			instruction.Bytes)
	}
	if instruction.Format == InstructionFormat1 {
		return fmt.Sprintf("%s : %s",
			instruction.InstructionAddress.StringHex(), instruction.Opcode.String())
	}
	if instruction.Format == InstructionFormat2 {
		return fmt.Sprintf("%s : %s\n    R1: %s\n    R2: %s",
			instruction.InstructionAddress.StringHex(), instruction.Opcode.String(),
			instruction.R1.String(), instruction.R2.String())
	}
	if instruction.IsFormatSIC34() {
		return fmt.Sprintf("%s : %s\n    Bytes: %X\n    Format: %s\n    AbsoluteAddressingMode: %s\n    RelativeAddressingMode: %s\n    Address: %s\n    Operand: %s",
			instruction.InstructionAddress.StringHex(), instruction.Opcode.String(),
			instruction.Bytes,
			instruction.Format.String(),
			instruction.AbsoluteAddressingMode.String(),
			instruction.RelativeAddressingMode.String(),
			instruction.Address.StringHex(), instruction.Operand.StringHex())
	}
	return "ERROR PRINTING INSTRUCTION\n"
}
