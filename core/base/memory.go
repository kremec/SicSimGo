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

	// TODO: Check if address range goes beyond memory size

	return units.Float48{
		memory.Data[address],
		memory.Data[address+1],
		memory.Data[address+2],
		memory.Data[address+3],
		memory.Data[address+4],
		memory.Data[address+5],
	}
}
func SetFloat(addressBytes units.Int24, value units.Float48) {
	address := toAddress(addressBytes)

	// TODO: Check if address range goes beyond memory size

	memory.Data[address] = value[0]
	memory.Data[address+1] = value[1]
	memory.Data[address+2] = value[2]
	memory.Data[address+3] = value[3]
	memory.Data[address+4] = value[4]
	memory.Data[address+5] = value[5]
}

func GetSlice(startAddress units.Int24, endAddress units.Int24) []byte {
	start := toAddress(startAddress)
	end := toAddress(endAddress)

	return memory.Data[start:end]
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
func String16Bytes(addressBytes units.Int24) string {
	address := toAddress(addressBytes)

	return fmt.Sprintf("%02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X",
		memory.Data[address],
		memory.Data[address+1],
		memory.Data[address+2],
		memory.Data[address+3],
		memory.Data[address+4],
		memory.Data[address+5],
		memory.Data[address+6],
		memory.Data[address+7],
		memory.Data[address+8],
		memory.Data[address+9],
		memory.Data[address+10],
		memory.Data[address+11],
		memory.Data[address+12],
		memory.Data[address+13],
		memory.Data[address+14],
		memory.Data[address+15],
	)
}
