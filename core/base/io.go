package base

import (
	"bufio"
	"fmt"
	"os"
)

/*
DEFINITIONS
*/
type Device byte

/*
DEBUG
*/
const debugRead bool = true
const debugWrite bool = true

/*
OPERATIONS
*/
func Test() bool {
	return true
}

func Read(device Device) (byte, error) {
	switch device {
	case Device(0x00):
		// Stdin
		if debugRead {
			fmt.Println("Reading from stdin")
		}
		reader := bufio.NewReader(os.Stdin)
		return readByte(reader)
	case Device(0x01):
		// Stdout
		if debugRead {
			fmt.Println("Reading from stdout")
		}
		reader := bufio.NewReader(os.Stdout)
		return readByte(reader)
	case Device(0x02):
		// Stderr
		if debugRead {
			fmt.Println("Reading from stderr")
		}
		reader := bufio.NewReader(os.Stderr)
		return readByte(reader)
	default:
		// XX.dev file
		filename := fmt.Sprintf("%02X.dev", device)
		if debugRead {
			fmt.Println("Reading from file ", filename)
		}
		file, err := os.Open(filename)
		if err != nil {
			return 0x00, err
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		return readByte(reader)
	}
}
func readByte(reader *bufio.Reader) (byte, error) {
	readByte, err := reader.ReadByte()
	if err != nil {
		panic("error reading byte")
	}
	if debugRead {
		fmt.Println("  Read byte:", readByte)
	}
	return readByte, err
}

func Write(device Device, data byte) error {

	if debugWrite {
		fmt.Println("Writing to device", device, "data:", data)
	}

	switch device {
	case Device(0x0):
		// Stdout
		_, err := os.Stdout.Write([]byte{data})
		return err
	case Device(0x1):
		// Stdout
		_, err := os.Stdout.Write([]byte{data})
		return err
	case Device(0x2):
		// Stderr
		_, err := os.Stderr.Write([]byte{data})
		return err
	}
	// XX.dev file
	filename := fmt.Sprintf("%02X.dev", device)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte{data})
	if err != nil {
		return err
	}

	return nil
}
