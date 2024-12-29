package assembly

import (
	"bufio"
	"fmt"
	"os"
	"sicsimgo/core/base"
	"sicsimgo/core/proc"
	"sicsimgo/core/units"
	"strconv"
	"strings"
)

/*
IMPLEMENTATIONS
*/
type Symbol struct {
	Name       string
	Address    units.Int24
	Data       bool
	DataLength int
}
type SymbolTable map[string]Symbol

var SyntaxNodes []SyntaxNode

/*
DEBUG
*/
var debugParseProgram bool = false

/*
OPERATIONS
*/
func LoadProgram(file *os.File) (string, units.Int24, map[units.Int24]proc.Instruction, SymbolTable) {
	var programName string
	var endPC units.Int24
	var disassembly map[units.Int24]proc.Instruction = make(map[units.Int24]proc.Instruction)
	var symbolTable SymbolTable = make(SymbolTable)

	// First pass
	LocationCounter := units.Int24{0x00, 0x00, 0x00}
	LineCounter := 0
	baseEnabled := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		LineCounter++
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		// Get syntax node
		syntaxNode := getSyntaxNode(line, LineCounter)
		if syntaxNode == nil {
			continue
		} else if syntaxNode.Mnemonic == "" && syntaxNode.Comment != "" {
			SyntaxNodes = append(SyntaxNodes, *syntaxNode)
			continue
		} else if syntaxNode.Mnemonic == "" {
			continue
		}

		syntaxNode.LocationCounter = LocationCounter

		// Directives
		switch syntaxNode.Mnemonic {
		case START:
			programName = syntaxNode.Label
		case ORG:
			orgValue := GetAbsoluteOperandAddress(syntaxNode.Operands[0])
			LocationCounter = orgValue
		case BASE:
			baseEnabled = true
		case NOBASE:
			baseEnabled = false
		case RESW:
			if syntaxNode.Label != "" {
				symbol := symbolTable[syntaxNode.Label]
				symbol.Name = syntaxNode.Label
				symbol.Address = LocationCounter
				symbol.Data = true
				symbol.DataLength = 3
				symbolTable[syntaxNode.Label] = symbol
			}
			LocationCounter = LocationCounter.Add(units.Int24{0x00, 0x00, 0x03})
		case WORD:
			if syntaxNode.Label != "" {
				symbol := symbolTable[syntaxNode.Label]
				symbol.Name = syntaxNode.Label
				symbol.Address = LocationCounter
				symbol.Data = true
				symbol.DataLength = 3
				symbolTable[syntaxNode.Label] = symbol
			}
			LocationCounter = LocationCounter.Add(units.Int24{0x00, 0x00, 0x03})
		case RESB:
			if syntaxNode.Label != "" {
				symbol := symbolTable[syntaxNode.Label]
				symbol.Name = syntaxNode.Label
				symbol.Address = LocationCounter
				symbol.Data = true
				symbol.DataLength = 1
				symbolTable[syntaxNode.Label] = symbol
			}
			LocationCounter = LocationCounter.Add(units.Int24{0x00, 0x00, 0x01})
		case BYTE:
			if syntaxNode.Label != "" {
				symbol := symbolTable[syntaxNode.Label]
				symbol.Name = syntaxNode.Label
				symbol.Address = LocationCounter
				symbol.Data = true
				symbol.DataLength = 1
				symbolTable[syntaxNode.Label] = symbol
			}
			LocationCounter = LocationCounter.Add(units.Int24{0x00, 0x00, 0x01})
		}

		// Instructions
		if IsMnemonicInstruction(syntaxNode.MnemonicType) {
			instruction := GetInstructionFromSyntaxNode(*syntaxNode, LocationCounter)
			if baseEnabled {
				instruction.RelativeAddressingMode = proc.BaseRelativeAddressing
			}
			disassembly[LocationCounter] = instruction
			if syntaxNode.Label != "" {
				symbol := symbolTable[syntaxNode.Label]
				symbol.Name = syntaxNode.Label
				symbol.Address = LocationCounter
				symbolTable[syntaxNode.Label] = symbol
			}

			switch instruction.Format {
			case proc.InstructionFormat1:
				LocationCounter = LocationCounter.Add(units.Int24{0x00, 0x00, 0x01})
			case proc.InstructionFormat2:
				LocationCounter = LocationCounter.Add(units.Int24{0x00, 0x00, 0x02})
			case proc.InstructionFormat3:
				LocationCounter = LocationCounter.Add(units.Int24{0x00, 0x00, 0x03})
			case proc.InstructionFormat4:
				LocationCounter = LocationCounter.Add(units.Int24{0x00, 0x00, 0x04})
			}
		}

		SyntaxNodes = append(SyntaxNodes, *syntaxNode)
	}

	if debugParseProgram {
		for _, syntaxNode := range SyntaxNodes {
			fmt.Println(syntaxNode.String())
			fmt.Println()
		}
		fmt.Println()
	}

	if debugParseProgram {
		for label, symbol := range symbolTable {
			fmt.Printf("Label: %s, Address: %s\n", label, symbol.Address.StringHex())
		}
		fmt.Println()
		fmt.Println()
	}

	// Second pass
	for _, syntaxNode := range SyntaxNodes {

		if syntaxNode.IsComment {
			continue
		}

		// Directives
		switch syntaxNode.Mnemonic {
		case EQU:
			equValue := GetOperandAddress(syntaxNode.Operands[0], symbolTable)
			if syntaxNode.Label != "" {
				symbol := symbolTable[syntaxNode.Label]
				symbol.Name = syntaxNode.Label
				symbol.Address = equValue
				symbol.Data = true
				symbol.DataLength = 3
				symbolTable[syntaxNode.Label] = symbol
			}
		case WORD:
			base.SetWord(syntaxNode.LocationCounter, GetAbsoluteOperandAddress(syntaxNode.Operands[0]))
		case BYTE:
			base.SetByte(syntaxNode.LocationCounter, GetAbsoluteOperandAddress(syntaxNode.Operands[0])[0])
		}

		// Instructions
		if IsMnemonicInstruction(syntaxNode.MnemonicType) {
			instruction := disassembly[syntaxNode.LocationCounter]

			switch syntaxNode.MnemonicType {
			case MnemonicF2N:
				// TODO: SYSCALL
			case MnemonicF2R:
				instruction.R1 = GetRegisterIdFromMnemonic(syntaxNode.Operands[0])
			case MnemonicF2RN:
				instruction.R1 = GetRegisterIdFromMnemonic(syntaxNode.Operands[0])
				n, _ := strconv.Atoi(syntaxNode.Operands[1])
				instruction.R2 = base.RegisterId(uint8(n - 1))
			case MnemonicF2RR:
				if strings.Contains(syntaxNode.Operands[0], ",") {
					r1r2Strings := strings.Split(syntaxNode.Operands[0], ",")
					instruction.R1 = GetRegisterIdFromMnemonic(r1r2Strings[0])
					instruction.R2 = GetRegisterIdFromMnemonic(r1r2Strings[1])
				} else {
					instruction.R1 = GetRegisterIdFromMnemonic(syntaxNode.Operands[0])
					instruction.R2 = GetRegisterIdFromMnemonic(syntaxNode.Operands[1])
				}
			case MnemonicF3:
				instruction.AbsoluteAddressingMode = proc.DirectAbsoluteAddressing
			case MnemonicF3M:
				pcAfterInstruction := syntaxNode.LocationCounter.Add(units.Int24{0x00, 0x00, 0x03})
				baseEnabled := instruction.RelativeAddressingMode == proc.BaseRelativeAddressing
				operandAddress, absoluteAddressingMode, relativeAddressingMode, indexAddressingMode := GetOperandAddressAddressingModes(syntaxNode.Operands[0], pcAfterInstruction, baseEnabled, symbolTable)
				instruction.Address = operandAddress
				instruction.AbsoluteAddressingMode = absoluteAddressingMode
				instruction.RelativeAddressingMode = relativeAddressingMode
				instruction.IndexAddressingMode = indexAddressingMode
			case MnemonicF4M:
				pcAfterInstruction := syntaxNode.LocationCounter.Add(units.Int24{0x00, 0x00, 0x04})
				operandAddress, absoluteAddressingMode, _, indexAddressingMode := GetOperandAddressAddressingModes(syntaxNode.Operands[0], pcAfterInstruction, false, symbolTable)
				instruction.Address = operandAddress
				instruction.RelativeAddressingMode = proc.DirectRelativeAddressing
				instruction.AbsoluteAddressingMode = absoluteAddressingMode
				instruction.IndexAddressingMode = indexAddressingMode
			}

			instruction.Bytes = instruction.GetInstructionBytes()

			instructionBytesAddress := instruction.InstructionAddress
			for i := 0; i < len(instruction.Bytes); i++ {
				base.SetByte(instructionBytesAddress, instruction.Bytes[i])
				instructionBytesAddress = instructionBytesAddress.Add(units.Int24{0x00, 0x00, 0x01})
			}

			disassembly[syntaxNode.LocationCounter] = instruction
		}
	}

	return programName, endPC, disassembly, symbolTable
}

func getSyntaxNode(line string, lineNumber int) *SyntaxNode {
	var syntaxNode SyntaxNode
	syntaxNode.LineNumber = lineNumber

	// Check for comment
	commentIndex := strings.Index(line, ".")
	if commentIndex != -1 {
		syntaxNode.Comment = strings.TrimSpace(line[commentIndex+1:])
		line = strings.TrimSpace(line[:commentIndex])

		if len(line) == 0 {
			syntaxNode.IsComment = true
			return &syntaxNode
		}
	}

	tokens := strings.Fields(line)

	// Check for label
	if mnemonicType := GetMnemonic(MnemonicName(tokens[0])); mnemonicType == MnemonicUnknown {
		// First token is label
		syntaxNode.Label = tokens[0]
		tokens = tokens[1:]
	}

	// Get mnemonic
	mnemonicType := GetMnemonic(MnemonicName(tokens[0]))
	if mnemonicType == MnemonicUnknown {
		ErrLabelWithoutMnemonic(syntaxNode.Label)
		return nil
	}
	syntaxNode.Mnemonic = MnemonicName(tokens[0])
	syntaxNode.MnemonicType = mnemonicType
	tokens = tokens[1:]

	// Get operands
	syntaxNode.Operands = tokens

	return &syntaxNode
}

func GetInstructionFromSyntaxNode(syntaxNode SyntaxNode, locationCounter units.Int24) proc.Instruction {
	instruction := proc.Instruction{}

	instruction.InstructionAddress = locationCounter

	switch syntaxNode.MnemonicType {
	case MnemonicF1:
		instruction.Format = proc.InstructionFormat1
	case MnemonicF2N, MnemonicF2R, MnemonicF2RN, MnemonicF2RR:
		instruction.Format = proc.InstructionFormat2
	case MnemonicF3, MnemonicF3M:
		instruction.Format = proc.InstructionFormat3
	case MnemonicF4M:
		instruction.Format = proc.InstructionFormat4
	}

	instruction.Opcode = GetInstructionOpcode(syntaxNode.Mnemonic)

	return instruction
}

func GetOperandAddressAddressingModes(operand string, pcFromLocationCounter units.Int24, baseEnabled bool, symbolTable SymbolTable) (units.Int24, proc.AbsoluteAddressingMode, proc.RelativeAddressingMode, proc.IndexAddressingMode) {
	// Absolute addressing mode
	var absoluteAddressingMode proc.AbsoluteAddressingMode
	if strings.HasPrefix(operand, "#") {
		absoluteAddressingMode = proc.ImmediateAbsoluteAddressing
		operand = operand[1:]
	} else if strings.HasPrefix(operand, "@") {
		absoluteAddressingMode = proc.IndirectAbsoluteAddressing
		operand = operand[1:]
	} else {
		absoluteAddressingMode = proc.DirectAbsoluteAddressing
	}

	// Index addressing mode
	var indexAddressingMode proc.IndexAddressingMode
	if strings.HasSuffix(operand, ",X") {
		indexAddressingMode = true
		operand = operand[:len(operand)-2]
	}

	// Operand address
	operandAddress := GetOperandAddress(operand, symbolTable)

	// Relative addressing
	var relativeAddressingMode proc.RelativeAddressingMode
	if baseEnabled {
		relativeAddressingMode = proc.BaseRelativeAddressing
	} else {
		switch absoluteAddressingMode {
		case proc.ImmediateAbsoluteAddressing: // signed absolute / pc / base
			// Try Direct-relative addressing
			if operandAddress.Compare(units.IntToInt24(-2048)) >= 0 && operandAddress.Compare(units.IntToInt24(2047)) < 0 {
				relativeAddressingMode = proc.DirectRelativeAddressing
			} else {
				// Try PC-relative addressing
				pcRelativeAddress := operandAddress.Sub(pcFromLocationCounter)
				if pcRelativeAddress.Compare(units.IntToInt24(-2048)) >= 0 && pcRelativeAddress.Compare(units.IntToInt24(2047)) < 0 {
					relativeAddressingMode = proc.PCRelativeAddressing
					operandAddress = pcRelativeAddress
				} else {
					relativeAddressingMode = proc.BaseRelativeAddressing
				}
			}
		case proc.IndirectAbsoluteAddressing: // pc, base, absolute
			// Try PC-relative addressing
			pcRelativeAddress := operandAddress.Sub(pcFromLocationCounter)
			if pcRelativeAddress.Compare(units.IntToInt24(-2048)) >= 0 && pcRelativeAddress.Compare(units.IntToInt24(2047)) < 0 {
				relativeAddressingMode = proc.PCRelativeAddressing
				operandAddress = pcRelativeAddress
			} else {
				relativeAddressingMode = proc.BaseRelativeAddressing
			}
		case proc.DirectAbsoluteAddressing: // pc / base / absolute / sic absolute
			// Try PC-relative addressing
			pcRelativeAddress := operandAddress.Sub(pcFromLocationCounter)
			if pcRelativeAddress.Compare(units.IntToInt24(-2048)) >= 0 && pcRelativeAddress.Compare(units.IntToInt24(2047)) < 0 {
				relativeAddressingMode = proc.PCRelativeAddressing
				operandAddress = pcRelativeAddress
			} else {
				relativeAddressingMode = proc.BaseRelativeAddressing
			}
		}
	}

	return operandAddress, absoluteAddressingMode, relativeAddressingMode, indexAddressingMode
}

func GetOperandAddress(operand string, symbolTable SymbolTable) units.Int24 {
	var operandAddress units.Int24
	// Check for relative symbol - label
	labelOperand, ok := symbolTable[operand]
	if ok {
		operandAddress = labelOperand.Address
	} else {
		// Check for absolute symbol
		operandAddress = GetAbsoluteOperandAddress(operand)
	}

	return operandAddress
}

func GetAbsoluteOperandAddress(operand string) units.Int24 {
	var operandAddress units.Int24
	// Check for absolute symbol - dec number
	intOperand, err := strconv.Atoi(operand)
	if err == nil {
		operandAddress = units.IntToInt24(intOperand)
	}
	// Check for absolute symbol - hex number
	if strings.HasPrefix(operand, "0x") {
		operand = operand[2:]
		intOperand, err := strconv.ParseUint(operand, 16, 8)
		if err == nil {
			operandAddress = units.IntToInt24(int(intOperand))
		}
	}
	//Check for absolute symbol - bin number
	if strings.HasPrefix(operand, "0b") {
		operand = operand[2:]
		intOperand, err := strconv.ParseUint(operand, 2, 8)
		if err == nil {
			operandAddress = units.IntToInt24(int(intOperand))
		}
	}

	return operandAddress
}
