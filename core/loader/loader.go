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

var ProgramName string
var StartPC units.Int24

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

	var loadedProgramType LoadedProgramType

	switch filepath.Ext(fileName) {
	case ".asm":
		loadedProgramType = Assembly
		ProgramName, StartPC, Disassembly, SymbolTable, SyntaxNodes = assembly.LoadProgram(file)
	case ".obj":
		loadedProgramType = Bytecode
		ProgramName, StartPC, Disassembly, LastInstructionByteAddress = bytecode.LoadProgram(file)
	default:
		return internal.DefaultWindowTitle, units.Int24{}, None
	}

	UpdateDisassemblyInstructionAddressOperands()
	UpdateInstructionList()
	UpdateSymbolTableList()

	return ProgramName, StartPC, loadedProgramType
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
		return
	}
	if !strings.HasSuffix(fileName, ".lst") {
		fileName += ".lst"
	}

	file, err := os.Create(fileName)
	if err != nil {
		return
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
	fileName, err := dialog.File().Filter("Object file", "obj").Title("Save object file").Save()
	if err != nil {
		return
	}
	if !strings.HasSuffix(fileName, ".obj") {
		fileName += ".obj"
	}

	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	// Generate text (T) & modification (M) sections
	var textSections []string
	var modificationSections []string
	var totalByteCount int

	var bytesBuffer []byte
	var lastByteAddress units.Int24 = InstructionList[0].InstructionAddress
	const maxBufferLength = 30
	for _, syntaxNode := range SyntaxNodes {

		goToNextTRecord := false
		// Determine bytes to add
		var bytesToAdd []byte
		if syntaxNode.Mnemonic == assembly.BYTE {
			absoluteOperandAddress := assembly.GetAbsoluteOperandAddress(syntaxNode.Operands[0])
			bytesToAdd = []byte{absoluteOperandAddress[0]}
		} else if syntaxNode.Mnemonic == assembly.WORD {
			absoluteOperandAddress := assembly.GetAbsoluteOperandAddress(syntaxNode.Operands[0])
			bytesToAdd = []byte{absoluteOperandAddress[0], absoluteOperandAddress[1], absoluteOperandAddress[2]}
		} else if syntaxNode.Mnemonic == assembly.RESB {
			goToNextTRecord = true
			lastByteAddress = lastByteAddress.Add(units.IntToInt24(1))
		} else if syntaxNode.Mnemonic == assembly.RESW {
			goToNextTRecord = true
			lastByteAddress = lastByteAddress.Add(units.IntToInt24(3))
		} else if assembly.IsMnemonicInstruction(syntaxNode.MnemonicType) {
			instruction := Disassembly[syntaxNode.LocationCounter]
			bytesToAdd = instruction.Bytes

			if instruction.RelativeAddressingMode == proc.DirectRelativeAddressing {
				modificationOffset := lastByteAddress.Add(units.IntToInt24(len(bytesBuffer) + 1))
				switch instruction.Format {
				case proc.InstructionFormatSIC | proc.InstructionFormat3:
					modificationRecord := fmt.Sprintf("M%X%02X\n", modificationOffset, 0x03)
					modificationSections = append(modificationSections, modificationRecord)
				case proc.InstructionFormat4:
					modificationRecord := fmt.Sprintf("M%X%02X\n", modificationOffset, 0x05)
					modificationSections = append(modificationSections, modificationRecord)
				}
			}
		}

		if goToNextTRecord {
			textRecord := fmt.Sprintf("T%X%02X%X\n", lastByteAddress, byte(len(bytesBuffer)), bytesBuffer)
			textSections = append(textSections, textRecord)
			totalByteCount += len(bytesBuffer)

			lastByteAddress = lastByteAddress.Add(units.IntToInt24(len(bytesBuffer)))
			bytesBuffer = []byte{}
		} else if len(bytesBuffer)+len(bytesToAdd) > maxBufferLength {
			appendBytes := bytesToAdd[:maxBufferLength-len(bytesBuffer)]
			leftOverBytes := bytesToAdd[maxBufferLength-len(bytesBuffer):]
			bytesBuffer = append(bytesBuffer, appendBytes...)

			textRecord := fmt.Sprintf("T%X%02X%X\n", lastByteAddress, byte(len(bytesBuffer)), bytesBuffer)
			textSections = append(textSections, textRecord)
			totalByteCount += len(bytesBuffer)

			lastByteAddress = lastByteAddress.Add(units.IntToInt24(len(bytesBuffer)))
			bytesBuffer = leftOverBytes
		} else {
			bytesBuffer = append(bytesBuffer, bytesToAdd...)
		}
	}
	if len(bytesBuffer) > 0 {
		textRecord := fmt.Sprintf("T%X%02X%X\n", lastByteAddress, byte(len(bytesBuffer)), bytesBuffer)
		textSections = append(textSections, textRecord)
		totalByteCount += len(bytesBuffer)
	}

	// Write header (H) section
	headerRecord := fmt.Sprintf("H%-6s%X%X\n", ProgramName, StartPC, units.IntToInt24(totalByteCount))
	file.WriteString(headerRecord)
	fmt.Println(headerRecord)

	// Write text (T) sections
	for _, textSection := range textSections {
		file.WriteString(textSection)
		fmt.Println(textSection)
	}

	// Write modification (M) sections
	for _, modificationSection := range modificationSections {
		file.WriteString(modificationSection)
		fmt.Println(modificationSection)
	}

	// Write end (E) section
	endRecord := fmt.Sprintf("E%X\n", LastInstructionByteAddress)
	file.WriteString(endRecord)
	fmt.Println(endRecord)
}
