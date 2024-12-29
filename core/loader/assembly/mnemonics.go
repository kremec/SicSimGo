package assembly

import (
	"sicsimgo/core/base"
	"sicsimgo/core/proc"
	"strings"
)

/*
IMPLEMENTATIONS
*/
type MnemonicName string

const (
	BASE   MnemonicName = "BASE"
	NOBASE MnemonicName = "NOBASE"
	LTORG  MnemonicName = "LTORG"
)

const (
	START MnemonicName = "START"
	END   MnemonicName = "END"
	ORG   MnemonicName = "ORG"
	EQU   MnemonicName = "EQU"
)

const (
	FIX   MnemonicName = "FIX"
	FLOAT MnemonicName = "FLOAT"
	HIO   MnemonicName = "HIO"
	NORM  MnemonicName = "NORM"
	SIO   MnemonicName = "SIO"
	TIO   MnemonicName = "TIO"
)

const (
	SVC MnemonicName = "SVC"
)

const (
	CLEAR MnemonicName = "CLEAR"
	TIXR  MnemonicName = "TIXR"
)

const (
	SHIFTL MnemonicName = "SHIFTL"
	SHIFTR MnemonicName = "SHIFTR"
)

const (
	ADDR  MnemonicName = "ADDR"
	COMPR MnemonicName = "COMPR"
	DIVR  MnemonicName = "DIVR"
	MULR  MnemonicName = "MULR"
	RMO   MnemonicName = "RMO"
	SUBR  MnemonicName = "SUBR"
)

const (
	RSUB MnemonicName = "RSUB"
)

const (
	ADD   MnemonicName = "ADD"
	ADDF  MnemonicName = "ADDF"
	AND   MnemonicName = "AND"
	COMP  MnemonicName = "COMP"
	COMPF MnemonicName = "COMPF"
	DIV   MnemonicName = "DIV"
	DIVF  MnemonicName = "DIVF"
	J     MnemonicName = "J"
	JEQ   MnemonicName = "JEQ"
	JGT   MnemonicName = "JGT"
	JLT   MnemonicName = "JLT"
	JSUB  MnemonicName = "JSUB"
	LDA   MnemonicName = "LDA"
	LDB   MnemonicName = "LDB"
	LDCH  MnemonicName = "LDCH"
	LDF   MnemonicName = "LDF"
	LDL   MnemonicName = "LDL"
	LDS   MnemonicName = "LDS"
	LDT   MnemonicName = "LDT"
	LDX   MnemonicName = "LDX"
	LPS   MnemonicName = "LPS"
	MUL   MnemonicName = "MUL"
	MULF  MnemonicName = "MULF"
	OR    MnemonicName = "OR"
	RD    MnemonicName = "RD"
	SSK   MnemonicName = "SSK"
	STA   MnemonicName = "STA"
	STB   MnemonicName = "STB"
	STCH  MnemonicName = "STCH"
	STF   MnemonicName = "STF"
	STI   MnemonicName = "STI"
	STL   MnemonicName = "STL"
	STS   MnemonicName = "STS"
	STSW  MnemonicName = "STSW"
	STT   MnemonicName = "STT"
	STX   MnemonicName = "STX"
	SUB   MnemonicName = "SUB"
	SUBF  MnemonicName = "SUBF"
	TD    MnemonicName = "TD"
	TIX   MnemonicName = "TIX"
	WD    MnemonicName = "WD"
)

const (
	BYTE MnemonicName = "BYTE"
	WORD MnemonicName = "WORD"
)

const (
	RESB MnemonicName = "RESB"
	RESW MnemonicName = "RESW"
)

type MnemonicType int

const (
	MnemonicDirective MnemonicType = iota
	MnemonicDirectiveN
	MnemonicF1
	MnemonicF2N
	MnemonicF2R
	MnemonicF2RN
	MnemonicF2RR
	MnemonicF3
	MnemonicF3M
	MnemonicF4M
	MnemonicStorageD
	MnemonicStorageN

	MnemonicUnknown
)

var Mnemonics = map[MnemonicName]MnemonicType{
	BASE:   MnemonicDirective,
	NOBASE: MnemonicDirective,
	LTORG:  MnemonicDirective,

	START: MnemonicDirectiveN,
	END:   MnemonicDirectiveN,
	ORG:   MnemonicDirectiveN,
	EQU:   MnemonicDirectiveN,

	FIX:   MnemonicF1,
	FLOAT: MnemonicF1,
	HIO:   MnemonicF1,
	NORM:  MnemonicF1,
	SIO:   MnemonicF1,
	TIO:   MnemonicF1,

	SVC: MnemonicF2N,

	CLEAR: MnemonicF2R,
	TIXR:  MnemonicF2R,

	SHIFTL: MnemonicF2RN,
	SHIFTR: MnemonicF2RN,

	ADDR:  MnemonicF2RR,
	COMPR: MnemonicF2RR,
	DIVR:  MnemonicF2RR,
	MULR:  MnemonicF2RR,
	RMO:   MnemonicF2RR,
	SUBR:  MnemonicF2RR,

	RSUB: MnemonicF3,

	ADD:   MnemonicF3M,
	ADDF:  MnemonicF3M,
	AND:   MnemonicF3M,
	COMP:  MnemonicF3M,
	COMPF: MnemonicF3M,
	DIV:   MnemonicF3M,
	DIVF:  MnemonicF3M,
	J:     MnemonicF3M,
	JEQ:   MnemonicF3M,
	JGT:   MnemonicF3M,
	JLT:   MnemonicF3M,
	JSUB:  MnemonicF3M,
	LDA:   MnemonicF3M,
	LDB:   MnemonicF3M,
	LDCH:  MnemonicF3M,
	LDF:   MnemonicF3M,
	LDL:   MnemonicF3M,
	LDS:   MnemonicF3M,
	LDT:   MnemonicF3M,
	LDX:   MnemonicF3M,
	LPS:   MnemonicF3M,
	MUL:   MnemonicF3M,
	MULF:  MnemonicF3M,
	OR:    MnemonicF3M,
	RD:    MnemonicF3M,
	SSK:   MnemonicF3M,
	STA:   MnemonicF3M,
	STB:   MnemonicF3M,
	STCH:  MnemonicF3M,
	STF:   MnemonicF3M,
	STI:   MnemonicF3M,
	STL:   MnemonicF3M,
	STS:   MnemonicF3M,
	STSW:  MnemonicF3M,
	STT:   MnemonicF3M,
	STX:   MnemonicF3M,
	SUB:   MnemonicF3M,
	SUBF:  MnemonicF3M,
	TD:    MnemonicF3M,
	TIX:   MnemonicF3M,
	WD:    MnemonicF3M,

	RESB: MnemonicStorageD,
	RESW: MnemonicStorageD,

	BYTE: MnemonicStorageN,
	WORD: MnemonicStorageN,
}

/*
OPERATIONS
*/
func GetMnemonic(mnemonicName MnemonicName) MnemonicType {
	mnemonicType, exists := Mnemonics[mnemonicName]
	if !exists {
		// Check for MnemonicF4M type
		if strings.HasPrefix(string(mnemonicName), "+") {
			mnemonicName = mnemonicName[1:]
			mnemonicType, exists = Mnemonics[mnemonicName]
			if exists && mnemonicType == MnemonicF3M {
				return MnemonicF4M
			}
		}
		return MnemonicUnknown
	}

	return mnemonicType
}

func IsMnemonicInstruction(mnemonicType MnemonicType) bool {
	switch mnemonicType {
	case MnemonicF1, MnemonicF2N, MnemonicF2R, MnemonicF2RN, MnemonicF2RR, MnemonicF3, MnemonicF3M, MnemonicF4M:
		return true
	}
	return false
}

func GetInstructionOpcode(mnemonic MnemonicName) proc.Opcode {
	switch mnemonic {
	case ADD:
		return proc.ADD
	case ADDF:
		return proc.ADDF
	case ADDR:
		return proc.ADDR
	case AND:
		return proc.AND
	case CLEAR:
		return proc.CLEAR
	case COMP:
		return proc.COMP
	case COMPF:
		return proc.COMPF
	case COMPR:
		return proc.COMPR
	case DIV:
		return proc.DIV
	case DIVF:
		return proc.DIVF
	case DIVR:
		return proc.DIVR
	case FIX:
		return proc.FIX
	case FLOAT:
		return proc.FLOAT
	case HIO:
		return proc.HIO
	case J:
		return proc.J
	case JEQ:
		return proc.JEQ
	case JGT:
		return proc.JGT
	case JLT:
		return proc.JLT
	case JSUB:
		return proc.JSUB
	case LDA:
		return proc.LDA
	case LDB:
		return proc.LDB
	case LDCH:
		return proc.LDCH
	case LDF:
		return proc.LDF
	case LDL:
		return proc.LDL
	case LDS:
		return proc.LDS
	case LDT:
		return proc.LDT
	case LDX:
		return proc.LDX
	case LPS:
		return proc.LPS
	case MUL:
		return proc.MUL
	case MULF:
		return proc.MULF
	case MULR:
		return proc.MULR
	case NORM:
		return proc.NORM
	case OR:
		return proc.OR
	case RD:
		return proc.RD
	case RMO:
		return proc.RMO
	case RSUB:
		return proc.RSUB
	case SHIFTL:
		return proc.SHIFTL
	case SHIFTR:
		return proc.SHIFTR
	case SIO:
		return proc.SIO
	case SSK:
		return proc.SSK
	case STA:
		return proc.STA
	case STB:
		return proc.STB
	case STCH:
		return proc.STCH
	case STF:
		return proc.STF
	case STI:
		return proc.STI
	case STL:
		return proc.STL
	case STS:
		return proc.STS
	case STSW:
		return proc.STSW
	case STT:
		return proc.STT
	case STX:
		return proc.STX
	case SUB:
		return proc.SUB
	case SUBF:
		return proc.SUBF
	case SUBR:
		return proc.SUBR
	case SVC:
		return proc.SVC
	case TD:
		return proc.TD
	case TIO:
		return proc.TIO
	case TIX:
		return proc.TIX
	case TIXR:
		return proc.TIXR
	case WD:
		return proc.WD
	}

	return proc.Opcode(0x00)
}

func GetRegisterIdFromMnemonic(operand string) base.RegisterId {
	switch operand {
	case "A":
		return base.RegisterAId
	case "X":
		return base.RegisterXId
	case "L":
		return base.RegisterLId
	case "B":
		return base.RegisterBId
	case "S":
		return base.RegisterSID
	case "T":
		return base.RegisterTId
	case "F":
		return base.RegisterFId
	case "PC":
		return base.RegisterPCId
	case "SW":
		return base.RegisterSWId
	}

	return base.RegisterId(0x00)
}
