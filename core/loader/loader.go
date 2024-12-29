package loader

import (
	"os"
	"path/filepath"
	"sicsimgo/core/loader/assembly"
	"sicsimgo/core/loader/bytecode"
	"sicsimgo/core/proc"
	"sicsimgo/core/units"
	"sicsimgo/internal"
	"sort"

	"github.com/sqweek/dialog"
)

/*
IMPLEMENTATION
*/
var Disassembly map[units.Int24]proc.Instruction
var InstructionList []proc.Instruction

var LastInstructionByteAddress units.Int24

var SymbolTable assembly.SymbolTable
var SymbolTableList []assembly.Symbol

/*
OPERATIONS
*/
func OpenAsmObjFile() (string, units.Int24) {
	filename, err := dialog.File().Filter("Assembly / Object files", "asm", "obj").Filter("Assembly files", "asm").Filter("Object files", "obj").Title("Select object / assembly file").Load()
	if err != nil {
		return internal.DefaultWindowTitle, units.Int24{}
	}

	file, err := os.Open(filename)
	if err != nil {
		return internal.DefaultWindowTitle, units.Int24{}
	}
	defer file.Close()

	var programName string
	var startPC units.Int24

	switch filepath.Ext(filename) {
	case ".asm":
		programName, startPC, Disassembly, SymbolTable = assembly.LoadProgram(file)
	case ".obj":
		programName, startPC, Disassembly, LastInstructionByteAddress = bytecode.LoadProgram(file)
	default:
		return internal.DefaultWindowTitle, units.Int24{}
	}

	UpdateDisassemblyInstructionAddressOperands()
	UpdateInstructionList()
	UpdateSymbolTableList()

	return programName, startPC
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
