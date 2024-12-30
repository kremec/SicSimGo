package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"sicsimgo/core/loader/assembly"
	"sicsimgo/core/loader/bytecode"
	"sicsimgo/core/proc"
	"sicsimgo/core/units"
	"sicsimgo/internal"
	"sort"
	"strings"

	"github.com/sqweek/dialog"
)

/*
DEFINITIONS
*/
type LoadedProgramType int

const (
	Assembly LoadedProgramType = iota
	Bytecode
	None
)

/*
IMPLEMENTATION
*/
var Disassembly map[units.Int24]proc.Instruction
var InstructionList []proc.Instruction

var LastInstructionByteAddress units.Int24

var SymbolTable assembly.SymbolTable
var SymbolTableList []assembly.Symbol

var SyntaxNodes []assembly.SyntaxNode

/*
OPERATIONS
*/
func OpenAsmObjFile() (string, units.Int24, LoadedProgramType) {
	fileName, err := dialog.File().Filter("Assembly / Object files", "asm", "obj").Filter("Assembly files", "asm").Filter("Object files", "obj").Title("Select object / assembly file").Load()
	if err != nil {
		return internal.DefaultWindowTitle, units.Int24{}, None
	}

	file, err := os.Open(fileName)
	if err != nil {
		return internal.DefaultWindowTitle, units.Int24{}, None
	}
	defer file.Close()

	var programName string
	var startPC units.Int24
	var loadedProgramType LoadedProgramType

	switch filepath.Ext(fileName) {
	case ".asm":
		loadedProgramType = Assembly
		programName, startPC, Disassembly, SymbolTable, SyntaxNodes = assembly.LoadProgram(file)
	case ".obj":
		loadedProgramType = Bytecode
		programName, startPC, Disassembly, LastInstructionByteAddress = bytecode.LoadProgram(file)
	default:
		return internal.DefaultWindowTitle, units.Int24{}, None
	}

	UpdateDisassemblyInstructionAddressOperands()
	UpdateInstructionList()
	UpdateSymbolTableList()

	return programName, startPC, loadedProgramType
}

func UpdateDisassemblyInstructionAddressOperands() {
	for address, instruction := range Disassembly {
		nextInstructionAddress := address
		for i := 0; i < len(instruction.Bytes); i++ {
			nextInstructionAddress = nextInstructionAddress.Add(units.Int24{0x00, 0x00, 0x01})
		}
		if instruction.IsFormatSIC34() {
			operand, address, relativeAddressingMode, indexAddressingMode, absoluteAddressingMode := instruction.GetOperandAddress(nextInstructionAddress)
			instruction.Operand = operand
			instruction.Address = address
			instruction.RelativeAddressingMode = relativeAddressingMode
			instruction.IndexAddressingMode = indexAddressingMode
			instruction.AbsoluteAddressingMode = absoluteAddressingMode
		}
		Disassembly[address] = instruction
	}
}

func UpdateInstructionList() {
	adresses := units.Int24Slice{}
	for key := range Disassembly {
		adresses = append(adresses, key)
	}
	sort.Sort(adresses)

	instructionList := make([]proc.Instruction, 0, len(Disassembly))
	for _, address := range adresses {
		instructionList = append(instructionList, Disassembly[address])
	}

	InstructionList = instructionList
}

func UpdateSymbolTableList() {
	for _, symbol := range SymbolTable {
		if symbol.Data {
			SymbolTableList = append(SymbolTableList, symbol)
		}
	}
	sort.Slice(SymbolTableList, func(i, j int) bool {
		return SymbolTableList[i].Address.Compare(SymbolTableList[j].Address) < 0
	})
}

func ResetDissasembly() {
	Disassembly = make(map[units.Int24]proc.Instruction)
	InstructionList = make([]proc.Instruction, 0)

	LastInstructionByteAddress = units.Int24{}

	SymbolTable = make(assembly.SymbolTable)
	SymbolTableList = make([]assembly.Symbol, 0)
}

func OutputLstFile() {
	fileName, err := dialog.File().Filter("List file", "lst").Title("Save list file").Save()
	if err != nil {
		panic("No file opened")
	}
	if !strings.HasSuffix(fileName, ".lst") {
		fileName += ".lst"
	}

	file, err := os.Create(fileName)
	if err != nil {
		panic("Cannot open file")
	}
	defer file.Close()

	currentLineNumber := 0
	for _, syntaxNode := range SyntaxNodes {
		for currentLineNumber < syntaxNode.LineNumber {
			file.WriteString("\n")
			currentLineNumber++
		}

		if syntaxNode.IsComment {
			file.WriteString(fmt.Sprintf("%-23s . %s\n", "", syntaxNode.Comment))
		} else if syntaxNode.MnemonicType == assembly.MnemonicDirective {
			if syntaxNode.Label != "" {
				file.WriteString(fmt.Sprintf("%s %s\n", syntaxNode.Label, syntaxNode.Mnemonic))
			} else {
				file.WriteString(fmt.Sprintf("%s\n", syntaxNode.Mnemonic))
			}
		} else if syntaxNode.MnemonicType == assembly.MnemonicDirectiveN {
			if syntaxNode.Label != "" {
				file.WriteString(fmt.Sprintf("%s %s %s\n", syntaxNode.Label, syntaxNode.Mnemonic, syntaxNode.Operands[0]))
			} else {
				file.WriteString(fmt.Sprintf("%s %s\n", syntaxNode.Mnemonic, syntaxNode.Operands[0]))
			}
		} else {
			file.WriteString(fmt.Sprintf("%-12s%-11X%-9s%-9s%s%-5s%s\n",
				syntaxNode.LocationCounter.StringHex(),
				func() []byte {
					if assembly.IsMnemonicInstruction(syntaxNode.MnemonicType) {
						return Disassembly[syntaxNode.LocationCounter].Bytes
					} else if syntaxNode.MnemonicType == assembly.MnemonicStorageD || syntaxNode.MnemonicType == assembly.MnemonicStorageN {
						return SymbolTable[syntaxNode.Label].Value
					}
					return []byte{}
				}(),
				syntaxNode.Label,
				syntaxNode.Mnemonic,
				strings.Join(syntaxNode.Operands, " "),
				"",
				func() string {
					if syntaxNode.Comment != "" {
						return ". " + syntaxNode.Comment
					}
					return ""
				}(),
			))
		}

		currentLineNumber++
	}
}

func OutputObjFile() {

}
