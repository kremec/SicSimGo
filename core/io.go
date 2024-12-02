package core

import (
	"bufio"
	"fmt"
	"os"
)

type Device byte

func Test() bool {
	return true
}

func Read(device Device) (byte, error) {
	switch device {
	case Device(0x0):
		// Stdin
		reader := bufio.NewReader(os.Stdin)
		return reader.ReadByte()
	case Device(0x1):
		// Stdout
		reader := bufio.NewReader(os.Stdout)
		return reader.ReadByte()
	case Device(0x2):
		// Stderr
		reader := bufio.NewReader(os.Stderr)
		return reader.ReadByte()
	}
	// XX.dev file
	filename := fmt.Sprintf("%02X.dev", device)
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	return reader.ReadByte()
}

func Write(device Device, data byte) error {
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
