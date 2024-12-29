package assembly

import (
	"fmt"
	"sicsimgo/core/units"
)

type SyntaxNode struct {
	Label        string
	Mnemonic     MnemonicName
	MnemonicType MnemonicType
	Operands     []string

	IsComment bool
	Comment   string

	LineNumber      int
	LocationCounter units.Int24
}

func (syntaxNode SyntaxNode) String() string {
	if syntaxNode.Mnemonic == "" && syntaxNode.Comment != "" {
		return fmt.Sprintf(".%s", syntaxNode.Comment)
	}
	return fmt.Sprintf("%s : %s\n    Operands: %s\n    Comment: %s", syntaxNode.LocationCounter.StringHex(), syntaxNode.Mnemonic, syntaxNode.Operands, syntaxNode.Comment)
}
