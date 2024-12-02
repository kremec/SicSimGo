package core

import (
	"fmt"
	"log"
)

func InvalidRegister(registerId RegisterId) error {
	log.Fatalf("Invalid register id: %d", registerId)
	return fmt.Errorf("Invalid register id: %d", registerId)
}
