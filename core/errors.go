package core

import (
	"fmt"
	"log"
)

func ErrNotImplemented(mnemonic string) error {
	log.Fatalf("Not implemented mnemonic: %s", mnemonic)
	return fmt.Errorf("Not implemented mnemonic: %s", mnemonic)
}

func ErrInvalidOpcode(opcode Opcode) error {
	return fmt.Errorf("Not implemented opcode: %02X", opcode)
}

func ErrInvalidAddressing() error {
	return fmt.Errorf("Invalid addressing with b=1 and p=1")
}

func ErrDisassemblyEmpty() error {
	return fmt.Errorf("Disassembly is empty")
}

func ErrDisassemblyIncorrect() error {
	log.Fatalf("Disassembly is incorrect")
	return fmt.Errorf("Disassembly is incorrect")
}
