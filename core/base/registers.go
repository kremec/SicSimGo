package base

import (
	"sicsimgo/core/units"
)

/*
DEFINITIONS
*/
type RegisterId uint8
type Registers struct {
	A  units.Int24
	X  units.Int24
	L  units.Int24
	B  units.Int24
	S  units.Int24
	T  units.Int24
	F  units.Float48
	PC units.Int24
	SW units.Int24 // 0x0, 0x40 in 0x80 represent "less than", "equal to" and "greater than" respectively
}

const (
	RegisterAId  RegisterId = 0
	RegisterXId  RegisterId = 1
	RegisterLId  RegisterId = 2
	RegisterBId  RegisterId = 3
	RegisterSID  RegisterId = 4
	RegisterTId  RegisterId = 5
	RegisterFId  RegisterId = 6
	RegisterPCId RegisterId = 8
	RegisterSWId RegisterId = 9
)

/*
IMPLEMENTATION
*/
var registers Registers = Registers{
	A:  units.Int24{},
	X:  units.Int24{},
	L:  units.Int24{},
	B:  units.Int24{},
	S:  units.Int24{},
	T:  units.Int24{},
	F:  units.Float48{},
	PC: units.Int24{},
	SW: units.Int24{},
}

/*
OPERATIONS
*/
func GetRegisterA() units.Int24 {
	return registers.A
}
func SetRegisterA(value units.Int24) {
	registers.A = value
}

func GetRegisterX() units.Int24 {
	return registers.X
}
func SetRegisterX(value units.Int24) {
	registers.X = value
}

func GetRegisterL() units.Int24 {
	return registers.L
}
func SetRegisterL(value units.Int24) {
	registers.L = value
}

func GetRegisterB() units.Int24 {
	return registers.B
}
func SetRegisterB(value units.Int24) {
	registers.B = value
}

func GetRegisterS() units.Int24 {
	return registers.S
}
func SetRegisterS(value units.Int24) {
	registers.S = value
}

func GetRegisterT() units.Int24 {
	return registers.T
}
func SetRegisterT(value units.Int24) {
	registers.T = value
}

func GetRegisterF() units.Float48 {
	return registers.F
}
func SetRegisterF(value units.Float48) {
	registers.F = value
}

func GetRegisterPC() units.Int24 {
	return registers.PC
}
func SetRegisterPC(value units.Int24) {
	registers.PC = value
}

func GetRegisterSW() units.Int24 {
	return registers.SW
}
func SetRegisterSW(value units.Int24) {
	registers.SW = value
}

func (registerId RegisterId) GetRegister() (units.Int24, error) {
	switch registerId {
	case RegisterAId:
		return registers.A, nil
	case RegisterXId:
		return registers.X, nil
	case RegisterLId:
		return registers.L, nil
	case RegisterBId:
		return registers.B, nil
	case RegisterSID:
		return registers.S, nil
	case RegisterTId:
		return registers.T, nil
	case RegisterFId:
		return units.Int24{registers.F[0], registers.F[1], registers.F[2]}, nil
	case RegisterPCId:
		return registers.PC, nil
	case RegisterSWId:
		return registers.SW, nil
	}

	return units.Int24{}, ErrInvalidRegister(registerId)
}
func (registerId RegisterId) SetRegister(value units.Int24) error {
	switch registerId {
	case RegisterAId:
		registers.A = value
	case RegisterXId:
		registers.X = value
	case RegisterLId:
		registers.L = value
	case RegisterBId:
		registers.B = value
	case RegisterSID:
		registers.S = value
	case RegisterTId:
		registers.T = value
	case RegisterFId:
		registers.F = units.Float48{value[0], value[1], value[2], 0x00, 0x00, 0x00}
	case RegisterPCId:
		registers.PC = value
	case RegisterSWId:
		registers.SW = value
	default:
		return ErrInvalidRegister(registerId)
	}
	return nil
}

func ResetRegisters() {
	resetValue := units.Int24{}
	registers.A = resetValue
	registers.X = resetValue
	registers.L = resetValue
	registers.B = resetValue
	registers.S = resetValue
	registers.T = resetValue
	registers.F = units.Float48{}
	registers.PC = resetValue
	registers.SW = resetValue
}

/*
STRINGS
*/
func (registerId RegisterId) String() string {
	switch registerId {
	case RegisterAId:
		return "A"
	case RegisterXId:
		return "X"
	case RegisterLId:
		return "L"
	case RegisterBId:
		return "B"
	case RegisterSID:
		return "S"
	case RegisterTId:
		return "T"
	case RegisterFId:
		return "F"
	case RegisterPCId:
		return "PC"
	case RegisterSWId:
		return "SW"
	}
	return "Not implemented"
}
