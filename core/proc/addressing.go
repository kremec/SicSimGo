package proc

import (
	"encoding/hex"
	"fmt"
	"sicsimgo/core/base"
	"sicsimgo/core/units"
)

/*
DEFINITIONS
*/
type RelativeAddressingMode int
type AbsoluteAddressingMode int
type IndexAddressingMode bool

const (
	DirectRelativeAddressing RelativeAddressingMode = 0
	PCRelativeAddressing     RelativeAddressingMode = 1
	BaseRelativeAddressing   RelativeAddressingMode = 2
	UnkownRelativeAddressing RelativeAddressingMode = 3
)

const (
	SICAbsoluteAddressing       AbsoluteAddressingMode = 0
	ImmediateAbsoluteAddressing AbsoluteAddressingMode = 1
	IndirectAbsoluteAddressing  AbsoluteAddressingMode = 2
	DirectAbsoluteAddressing    AbsoluteAddressingMode = 3
)

/*
DEBUG
*/
const debugGetOperandAddress bool = false

/*
OPERATIONS
*/
func GetRelativeAdressingModes(b, p bool) (RelativeAddressingMode, error) {
	if !b && !p {
		return DirectRelativeAddressing, nil
	} else if !b && p {
		return PCRelativeAddressing, nil
	} else if b && !p {
		return BaseRelativeAddressing, nil
	} else {
		return UnkownRelativeAddressing, ErrInvalidAddressing()
	}
}

func GetAbsoluteAdressingModes(n, i bool) AbsoluteAddressingMode {
	if !n && !i {
		return SICAbsoluteAddressing
	} else if !n && i {
		return ImmediateAbsoluteAddressing
	} else if n && !i {
		return IndirectAbsoluteAddressing
	} else {
		return DirectAbsoluteAddressing
	}
}

func GetIndexAdressingModes(x bool) IndexAddressingMode {
	if x {
		return true
	} else {
		return false
	}
}

func (instruction Instruction) GetNIXBPEBits() (n, i, x, b, p, e bool) {
	if instruction.IsFormatSIC34() && len(instruction.Bytes) > 2 {
		n = instruction.Bytes[0]&0b00000010 > 0
		i = instruction.Bytes[0]&0b00000001 > 0
		x = instruction.Bytes[1]&0b10000000 > 0
		b = instruction.Bytes[1]&0b01000000 > 0
		p = instruction.Bytes[1]&0b00100000 > 0
		e = instruction.Bytes[1]&0b00010000 > 0
	} else {
		n = false
		i = false
		x = false
		b = false
		p = false
		e = false
	}

	return n, i, x, b, p, e
}

func (instruction Instruction) GenerateNIXBPEBits() (n, i, x, b, p, e bool) {
	switch instruction.AbsoluteAddressingMode {
	case SICAbsoluteAddressing:
		n, i = false, false
	case ImmediateAbsoluteAddressing:
		n, i = false, true
	case IndirectAbsoluteAddressing:
		n, i = true, false
	case DirectAbsoluteAddressing:
		n, i = true, true
	}

	if instruction.IndexAddressingMode {
		x = true
	}

	switch instruction.RelativeAddressingMode {
	case DirectRelativeAddressing:
		b, p = false, false
	case PCRelativeAddressing:
		b, p = false, true
	case BaseRelativeAddressing:
		b, p = true, false
	case UnkownRelativeAddressing:
		b, p = true, true
	}

	if instruction.Format == InstructionFormat4 {
		e = true
	}

	return n, i, x, b, p, e
}

func (instruction Instruction) GetOperandAddress(pc units.Int24) (units.Int24, units.Int24, RelativeAddressingMode, IndexAddressingMode, AbsoluteAddressingMode) {

	if instruction.Format != InstructionFormatSIC && instruction.Format != InstructionFormat3 && instruction.Format != InstructionFormat4 {
		// Invalid instruction format
		return units.Int24{}, units.Int24{}, UnkownRelativeAddressing, IndexAddressingMode(false), SICAbsoluteAddressing
	}

	var address units.Int24
	var operand units.Int24

	// Get instruction addressing modes
	n, i, x, b, p, _ := instruction.GetNIXBPEBits()
	relativeAddressingMode, err := GetRelativeAdressingModes(b, p)
	if err != nil {
		// Invalid relative addressing
		return units.Int24{}, units.Int24{}, UnkownRelativeAddressing, IndexAddressingMode(false), SICAbsoluteAddressing
	}
	absoluteAddressingMode := GetAbsoluteAdressingModes(n, i)
	indexAddressingMode := GetIndexAdressingModes(x)

	// Get instruction operand address
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
		fmt.Println("    Instruction bytes: ", hex.EncodeToString(instruction.Bytes))
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
		address = address.Add(base.GetRegisterB())
	}

	if debugGetOperandAddress {
		fmt.Println("    Address after relative addressing: ", address.StringHex())
		fmt.Println("        PC: ", pc.StringHex())
	}

	// Index addressing
	if indexAddressingMode {
		address = address.Add(base.GetRegisterX())
	}

	if debugGetOperandAddress {
		fmt.Println("    Address after index addressing: ", address.StringHex())
	}

	// Absolute addressing
	// If instruction is a jump instruction, the address is the destination address, no operand needed
	if !instruction.IsJumpInstruction() {
		switch absoluteAddressingMode {
		case SICAbsoluteAddressing:
			operand = base.GetWord(address)
		case ImmediateAbsoluteAddressing:
			operand = address
		case IndirectAbsoluteAddressing:
			operand = base.GetWord(base.GetWord(address))
		case DirectAbsoluteAddressing:
			operand = base.GetWord(address)
		}

		if debugGetOperandAddress {
			fmt.Println("    Operand:", operand.StringHex())
			fmt.Println()
		}
	}

	return operand, address, relativeAddressingMode, indexAddressingMode, absoluteAddressingMode
}

func GetR1R2FromByte(byte2 byte) (base.RegisterId, base.RegisterId) {
	r1Id := base.RegisterId(byte2 & 0xF0 >> 4)
	r2Id := base.RegisterId(byte2 & 0x0F)
	return r1Id, r2Id
}

/*
STRINGS
*/
func (relativeAddressingMode RelativeAddressingMode) String() string {
	switch relativeAddressingMode {
	case DirectRelativeAddressing:
		return "Direct"
	case PCRelativeAddressing:
		return "PC-relative"
	case BaseRelativeAddressing:
		return "Base-relative"
	case UnkownRelativeAddressing:
		return "Unknown"
	}
	return "Not implemented"
}

func (absoluteAddressingMode AbsoluteAddressingMode) String() string {
	switch absoluteAddressingMode {
	case SICAbsoluteAddressing:
		return "SIC"
	case ImmediateAbsoluteAddressing:
		return "Immediate"
	case IndirectAbsoluteAddressing:
		return "Indirect"
	case DirectAbsoluteAddressing:
		return "Direct"
	}
	return "Not implemented"
}
