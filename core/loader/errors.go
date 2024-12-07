package loader

import (
	"fmt"
	"log"
)

func ErrDisassemblyEmpty() error {
	return fmt.Errorf("Disassembly is empty")
}

func ErrDisassemblyIncorrect() error {
	log.Fatalf("Disassembly is incorrect")
	return fmt.Errorf("Disassembly is incorrect")
}
