package loader

import (
	"bufio"
	"fmt"
	"os"
	"sicsimgo/core"
	"sicsimgo/core/units"
	"sicsimgo/internal"
	"strconv"
	"strings"

	"github.com/sqweek/dialog"
)

func OpenObjectFile() string {
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

var debugLoadProgram bool = true

func LoadProgram(file *os.File) string {
	var programName string
	var codeOffset units.Int24

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := scanner.Text()
		if debugLoadProgram {
			fmt.Println(record)
		}
		if record[0] == 'H' {
			progName, codeAddr, codeLen := GetHeaderRecord(record)
			if debugLoadProgram {
				fmt.Printf("Header: %s|%s|%s\n", progName, codeAddr.StringHex(), codeLen.StringHex())
			}

			programName = progName
			codeOffset = codeAddr
		} else if record[0] == 'T' {
			codeAddress, code := GetTextRecord(record)
			if debugLoadProgram {
				fmt.Printf("    Text: %s|% X\n", codeAddress.StringHex(), code)
			}

			idx := units.Int24{}
			for i := 0; i < len(code); i++ {
				core.SetByte(codeAddress.Add(idx.Add(codeOffset)), code[i])
				idx = idx.Add(units.Int24{0x00, 0x00, 0x01})
			}
		} else if record[0] == 'E' {
			endAddress := GetEndRecord(record)
			if debugLoadProgram {
				fmt.Printf("    End: %s\n", endAddress.StringHex())
			}

			core.SetRegisterPC(endAddress.Add(codeOffset))
		}
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