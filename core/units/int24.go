package units

import (
	"fmt"
	"strconv"
)

const (
	WORD_SIZE int = 3
)

type Int24 [3]byte
type Int24Slice []Int24

func (i Int24) ToUint32() uint32 {
	return uint32(i[0])<<16 | uint32(i[1])<<8 | uint32(i[2])
}

func (i Int24) IsNegative() bool {
	return i[0]&0b10000000 != 0
}
func (i Int24) ToInt32() int32 {
	if i.IsNegative() {
		return int32(-1)<<24 | int32(i[0])<<16 | int32(i[1])<<8 | int32(i[2])
	}
	return int32(i[0])<<16 | int32(i[1])<<8 | int32(i[2])
}

func ToInt24(s string) Int24 {
	if len(s) != WORD_SIZE*2 {
		return Int24{}
	}

	var result Int24
	for i := 0; i < 3; i++ {
		hexByte := s[i*2 : i*2+2]
		parsedByte, err := strconv.ParseUint(hexByte, 16, 8)
		if err != nil {
			panic(err)
		}
		result[i] = byte(parsedByte)
	}

	return result
}

/*
ARITHMETIC OPERATORS
*/
func (i Int24) Add(other Int24) Int24 {
	var result Int24
	var carry byte

	// Add first byte (LSB)
	sum := i[2] + other[2]
	carry = (i[2]&other[2] | (i[2]|other[2])&^sum) >> 0x7
	result[2] = sum

	// Add second byte with carry
	sum = i[1] + other[1] + carry
	carry = (i[1]&other[1] | (i[1]|other[1])&^sum) >> 7
	result[1] = sum

	// Add third byte (MSB) with carry and mask
	sum = i[0] + other[0] + carry
	result[0] = sum & 0x7F // Mask to 24 bits

	return result
}

func (i Int24) Sub(other Int24) Int24 {
	var result Int24
	var borrow byte

	// Subtract first byte (LSB)
	diff := i[2] - other[2]
	borrow = (^i[2]&other[2] | (^i[2]|other[2])&diff) >> 7
	result[2] = diff

	// Subtract second byte with borrow
	diff = i[1] - other[1] - borrow
	borrow = (^i[1]&other[1] | (^i[1]|other[1])&diff) >> 7
	result[1] = diff

	// Subtract third byte (MSB) with borrow and mask
	diff = i[0] - other[0] - borrow
	result[0] = diff & 0x7F // Mask to 24 bits

	return result
}

// Does signed multiplication, masks the result to 24 bits
func (i Int24) Mul(other Int24) Int24 {
	// Convert Int24 to int32 for multiplication
	a := i.ToInt32()
	b := other.ToInt32()

	// Perform multiplication
	product := a * b

	// Mask to fit into 24 bits, preserving two's complement
	product &= 0xFFFFFF // Keep only the lower 24 bits

	// Create new Int24 result
	var result Int24
	result[0] = byte((product >> 16) & 0xFF)
	result[1] = byte((product >> 8) & 0xFF)
	result[2] = byte(product & 0xFF)

	return result
}

// Does signed division, masks the result to 24 bits
func (i Int24) Div(other Int24) Int24 {
	// Convert Int24 to int32 for division
	a := i.ToInt32()
	b := other.ToInt32()

	// Perform division
	quotient := a / b

	// Mask to fit into 24 bits, preserving two's complement
	quotient &= 0xFFFFFF // Keep only the lower 24 bits

	// Create new Int24 result
	var result Int24
	result[0] = byte((quotient >> 16) & 0xFF)
	result[1] = byte((quotient >> 8) & 0xFF)
	result[2] = byte(quotient & 0xFF)

	return result
}

/*
BITWISE LOGICAL OPERATORS
*/
func (i Int24) And(other Int24) Int24 {
	var result Int24
	result[0] = i[0] & other[0]
	result[1] = i[1] & other[1]
	result[2] = i[2] & other[2]
	return result
}

func (i Int24) Or(other Int24) Int24 {
	var result Int24
	result[0] = i[0] | other[0]
	result[1] = i[1] | other[1]
	result[2] = i[2] | other[2]
	return result
}

func (i Int24) Xor(other Int24) Int24 {
	var result Int24
	result[0] = i[0] ^ other[0]
	result[1] = i[1] ^ other[1]
	result[2] = i[2] ^ other[2]
	return result
}

func (i Int24) Not() Int24 {
	var result Int24
	result[0] = ^i[0]
	result[1] = ^i[1]
	result[2] = ^i[2]
	return result
}

/*
LOGICAL OPERATORS
*/
func (i Int24) Compare(other Int24) int {
	for idx := 2; idx >= 0; idx-- {
		if i[idx] < other[idx] {
			return -1 // i < other
		}
		if i[idx] > other[idx] {
			return 1 // i > other
		}
	}
	return 0 // i == other
}
func (s Int24Slice) Len() int {
	return len(s)
}
func (s Int24Slice) Less(i, j int) bool {
	return s[i].Compare(s[j]) == -1
}
func (s Int24Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

/*
STRING
*/
func (i Int24) StringDecUnsigned() string {
	return fmt.Sprintf("%d", i.ToUint32())
}
func (i Int24) StringDecSigned() string {
	return fmt.Sprintf("%d", i.ToInt32())
}
func (i Int24) StringHex() string {
	return fmt.Sprintf("%02X %02X %02X", i[0], i[1], i[2])
}
func (i Int24) StringBin() string {
	return fmt.Sprintf("%08b %08b %08b", i[0], i[1], i[2])
}
