package proc

import (
	"fmt"
)

func ErrInvalidOpcode(opcode Opcode) error {
	return fmt.Errorf("Not implemented opcode: %02X", opcode)
}

func ErrInvalidAddressing() error {
	return fmt.Errorf("Invalid addressing with b=1 and p=1")
}
