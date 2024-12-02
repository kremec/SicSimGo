package proc

type Opcode byte

var (
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

type InstructionFormat int

const (
	InstructionFormat1   InstructionFormat = 1
	InstructionFormat2   InstructionFormat = 2
	InstructionFormatSIC InstructionFormat = 0
	InstructionFormat3   InstructionFormat = 3
	InstructionFormat4   InstructionFormat = 4

	InstructionFormat34 InstructionFormat = 34
	InstructionUnknown  InstructionFormat = -1
)

func GetInstructionFormatFromOpcode(opcode Opcode, byte1 byte) (InstructionFormat, error) {
	// Format 1
	format1Instructions := []Opcode{
		FIX, FLOAT, HIO, NORM, SIO, TIO,
	}
	for _, op := range format1Instructions {
		if opcode == op {
			return InstructionFormat1, nil
		}
	}

	// Format 2
	format2Instructions := []Opcode{
		ADDR, CLEAR, COMPR, DIVR, MULR, RMO, SHIFTL, SHIFTR, SUBR, SVC, TIXR,
	}
	for _, op := range format2Instructions {
		if opcode == op {
			return InstructionFormat2, nil
		}
	}

	// SIC format
	if (byte1 & 0b00000011) == 0x00 {
		return InstructionFormatSIC, nil
	}

	// Default to Format 3/4
	format34Instructions := []Opcode{
		ADD, ADDF, AND, COMP, COMPF, DIV, DIVF, J, JEQ, JGT, JLT, JSUB, LDA, LDB, LDCH, LDF, LDL, LDS, LDT, LDX, LPS, MUL, MULF, OR, RD, RSUB, SSK, STA, STB, STCH, STF, STI, STL, STS, STSW, STT, STX, SUB, SUBF, TD, TIX, WD,
	}
	for _, op := range format34Instructions {
		if opcode == op {
			return InstructionFormat34, nil
		}
	}

	return InstructionUnknown, InvalidOpcode(opcode)
}

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
