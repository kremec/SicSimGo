package units

import "fmt"

type Float48 [6]byte

/*
STRING
*/
func (f Float48) StringDec() string {
	return fmt.Sprintf("TODO") // TODO
}
func (f Float48) StringHex() string {
	return fmt.Sprintf("%02X %02X %02X %02X %02X %02X", f[0], f[1], f[2], f[3], f[4], f[5])
}
func (f Float48) StringBin() string {
	return fmt.Sprintf("%08b %08b %08b %08b %08b %08b", f[0], f[1], f[2], f[3], f[4], f[5])
}
