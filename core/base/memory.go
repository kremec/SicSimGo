package base

import (
	"fmt"
	"sicsimgo/core/units"
)

/*
DEFINITIONS
*/
type Memory struct {
	Data []byte
}

const (
	MEMORY_SIZE uint32 = 0x100000
	MAX_ADDRESS uint32 = 0xFFFFF
)

/*
IMPLEMENTATION
*/
var memory Memory = Memory{
	Data: make([]byte, MEMORY_SIZE),
}

/*
TRANSFORMATIONS
*/
func ToAddress(val uint32) units.Int24 {

	if val > MAX_ADDRESS {
		panic("Address out of range")
	}

	return units.Int24{
		byte((val >> 16) & 0xFF),
		byte((val >> 8) & 0xFF),
		byte(val & 0xFF),
	}
}

func toAddress(val units.Int24) uint32 {
	return val.ToUint32()
}

/*
OPERATIONS
*/
func GetByte(addressBytes units.Int24) byte {
	address := toAddress(addressBytes)

	return memory.Data[address]
}
func SetByte(addressBytes units.Int24, value byte) {
	address := toAddress(addressBytes)
	memory.Data[address] = value
}

func GetWord(addressBytes units.Int24) units.Int24 {
	address := toAddress(addressBytes)

	return units.Int24{
		memory.Data[address],
		memory.Data[address+1],
		memory.Data[address+2],
	}
}
func SetWord(addressBytes units.Int24, value units.Int24) {
	address := toAddress(addressBytes)
	memory.Data[address] = value[0]
	memory.Data[address+1] = value[1]
	memory.Data[address+2] = value[2]
}

func GetFloat(addressBytes units.Int24) units.Float48 {
	address := toAddress(addressBytes)

	float := units.Float48{}
	for i := 0; i < 6; i++ {
		byteAddress := address + uint32(i)
		if byteAddress > MAX_ADDRESS {
			panic("Address out of range")
		}
		float[i] = memory.Data[byteAddress]
	}

	return float
}
func SetFloat(addressBytes units.Int24, value units.Float48) {
	address := toAddress(addressBytes)

	for i := 0; i < 6; i++ {
		byteAddress := address + uint32(i)
		if byteAddress > MAX_ADDRESS {
			panic("Address out of range")
		}
		memory.Data[byteAddress] = value[i]
	}
}

func GetSlice(startAddress units.Int24, endAddress units.Int24) []byte {
	start := toAddress(startAddress)
	end := toAddress(endAddress)

	return memory.Data[start:end]
}

func GetSlice16(startAddress units.Int24) []byte {
	start := toAddress(startAddress)

	return []byte{
		memory.Data[start],
		memory.Data[start+1],
		memory.Data[start+2],
		memory.Data[start+3],
		memory.Data[start+4],
		memory.Data[start+5],
		memory.Data[start+6],
		memory.Data[start+7],
		memory.Data[start+8],
		memory.Data[start+9],
		memory.Data[start+10],
		memory.Data[start+11],
		memory.Data[start+12],
		memory.Data[start+13],
		memory.Data[start+14],
		memory.Data[start+15],
	}
}

func ResetMemory() {
	memory.Data = make([]byte, MEMORY_SIZE)
}

/*
STRINGS
*/
func StringAddress(addressBytes units.Int24) string {
	address := toAddress(addressBytes)

	return fmt.Sprintf("%05X", address&0xFFFFF)
}
