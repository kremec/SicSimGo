package loader

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"sicsimgo/core/base"
	"sicsimgo/core/units"
	"sicsimgo/internal"

	"github.com/sqweek/dialog"
)

/*
IMPLEMENTATIONS
*/
var LastInstructionByteAddress units.Int24

/*
DEBUG
*/
var debugLoadProgram bool = false

/*
OPERATIONS
*/
func OpenLoadObjFile() string {
	filename, err := dialog.File().Filter("Object files", "obj").Title("Select object file").Load()
	if err != nil {
		return internal.DefaultWindowTitle
	}

	file, err := os.Open(filename)
	if err != nil {
		return internal.DefaultWindowTitle
	}
	defer file.Close()

	return LoadProgram(file)
}

func LoadProgram(file *os.File) string {
	var programName string
	var codeOffset units.Int24

	var previousTextRecordAddr units.Int24
	var previousTextrecordCodeLen int
	var leftoverBytes []byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := scanner.Text()
		if debugLoadProgram {
			fmt.Println(record)
		}
		if record[0] == 'H' {
			progName, codeAddr, codeLen := GetHeaderRecord(record)
			if debugLoadProgram {
				fmt.Printf("  Header: %s|%s|%s\n", progName, codeAddr.StringHex(), codeLen.StringHex())
			}

			programName = progName
			codeOffset = codeAddr
			leftoverBytes = []byte{}
		} else if record[0] == 'T' {
			codeAddress, code := GetTextRecord(record)
			if debugLoadProgram {
				fmt.Printf("  Text: %s|% X\n", codeAddress.StringHex(), code)
			}

			// Update last instruction byte address
			LastInstructionByteAddress = codeAddress
			if debugLoadProgram {
				fmt.Printf("Fixing last instruction byte address - %s: adding %d\n", LastInstructionByteAddress.StringHex(), len(code))
			}
			for i := 0; i < len(code)-1; i++ {
				LastInstructionByteAddress = LastInstructionByteAddress.Add(units.Int24{0x00, 0x00, 0x01})
			}

			// Memory
			idx := units.Int24{}
			for i := 0; i < len(code); i++ {
				base.SetByte(codeAddress.Add(idx.Add(codeOffset)), code[i])
				idx = idx.Add(units.Int24{0x00, 0x00, 0x01})
			}

			// Handle leftover bytes
			if len(leftoverBytes) > 0 {
				if debugLoadProgram {
					fmt.Printf("Leftover bytes: % X\n", leftoverBytes)
				}
				savePreviousTextRecordAddr := previousTextRecordAddr
				// Check for text record continuity if leftover bytes exist
				for i := 0; i < previousTextrecordCodeLen; i++ {
					previousTextRecordAddr = previousTextRecordAddr.Add(units.Int24{0x00, 0x00, 0x01})
				}
				if codeAddress.Compare(previousTextRecordAddr) != 0 {
					if debugLoadProgram {
						fmt.Printf("codeAddress: %s, previousTextRecordAddr: %s, previousTextrecordCodeLen: %d\n", codeAddress.StringHex(), savePreviousTextRecordAddr.StringHex(), previousTextrecordCodeLen)
					}
					// Text records aren't continuing, error
					panic(fmt.Errorf("Previous T record ended with leftover bytes, but current T record isn't continuing from there onwards!"))
				}

				// Add leftover bytes
				newCode := leftoverBytes
				newCode = append(newCode, code...)
				code = newCode

				// Fix start address
				for i := 0; i < len(leftoverBytes); i++ {
					codeAddress = codeAddress.Sub(units.Int24{0x00, 0x00, 0x01})
				}

				leftoverBytes = []byte{}
			}

			// Dissasembly
			instructions, bytesFromIncompleteInstruction := GetInstructions(codeAddress, code)
			for address, instruction := range instructions {
				if debugLoadProgram {
					fmt.Printf("    Address: %s, Bytes: % X, Format: %s, Opcode: %s, Operand: %s\n", address.StringHex(), instruction.Bytes, instruction.Format.String(), instruction.Opcode.String(), instruction.Operand.StringHex())
				}
				Disassembly[address] = instruction
			}

			// Prepare next text record for leftover bytes
			if len(bytesFromIncompleteInstruction) != 0 {
				leftoverBytes = bytesFromIncompleteInstruction
				previousTextRecordAddr = codeAddress
				previousTextrecordCodeLen = len(code)
			}
		} else if record[0] == 'E' {
			endAddress := GetEndRecord(record)
			if debugLoadProgram {
				fmt.Printf("  End: %s\n", endAddress.StringHex())
			}

			base.SetRegisterPC(endAddress.Add(codeOffset))
			leftoverBytes = []byte{}
		}
	}

	UpdateDisassemblyInstructionAdressOperands()
	UpdateDisassemblyInstructionList()

	if debugLoadProgram {
		fmt.Printf("Last instruction byte address: %s\n", LastInstructionByteAddress.StringHex())
	}
	return programName
}

func GetHeaderRecord(record string) (string, units.Int24, units.Int24) {
	programName := record[1:7]
	codeAddressStr := record[7:13]
	codeLengthStr := record[13:19]

	return strings.TrimSpace(programName), units.ToInt24(codeAddressStr), units.ToInt24(codeLengthStr)
}

func GetTextRecord(record string) (units.Int24, []byte) {
	codeAddressStr := record[1:7]
	codeLenStr := record[7:9]
	var code []byte

	codeLen, _ := strconv.ParseUint(codeLenStr, 16, 8)
	for i := 0; i < int(codeLen); i++ {
		b, _ := strconv.ParseUint(record[9+i*2:9+i*2+2], 16, 8)
		code = append(code, byte(b))
	}

	return units.ToInt24(codeAddressStr), code
}

func GetEndRecord(record string) units.Int24 {
	startAddressStr := record[1:7]

	return units.ToInt24(startAddressStr)
}
