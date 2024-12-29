package assembly

import (
	"fmt"
)

func ErrLabelWithoutMnemonic(label string) error {
	return fmt.Errorf("Label without mnemonic: %s", label)
}
