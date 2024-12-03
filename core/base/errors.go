package base

import (
	"fmt"
	"log"
)

func ErrInvalidRegister(registerId RegisterId) error {
	log.Fatalf("Invalid register id: %d", registerId)
	return fmt.Errorf("Invalid register id: %d", registerId)
}
