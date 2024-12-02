package proc

import (
	"fmt"
	"log"
)

func NotImplemented(mnemonic string) error {
	log.Fatalf("Not implemented mnemonic: %s", mnemonic)
	return fmt.Errorf("Not implemented mnemonic: %s", mnemonic)
}

func InvalidOpcode(opcode Opcode) error {
	log.Fatalf("Not implemented opcode: %02X", opcode)
	return fmt.Errorf("Not implemented opcode: %02X", opcode)
}

func InvalidAddressing() error {
	log.Fatalf("Invalid addressing with b=1 and p=1")
	return fmt.Errorf("Invalid addressing with b=1 and p=1")
}
