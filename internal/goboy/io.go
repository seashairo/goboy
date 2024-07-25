package goboy

import (
	"encoding/hex"
	"fmt"
	"os"
)

type IO struct {
	interrupts InterruptRegister
}

func NewIO() IO {
	return IO{
		interrupts: NewInterruptRegister(0),
	}
}

func (io *IO) writeByte(address uint16, value byte) {
	if address == 0xFF01 {
		appendSerialToFile(value)
		return
	}

	if address == 0xFF0F {
		io.interrupts.writeByte(value)
		return
	}

	fmt.Printf("Writing to %2.2X not supported (IO_REGISTERS)\n", address)
}

func (io *IO) readByte(address uint16) byte {
	if address == 0xFF0F {
		return io.interrupts.readByte()
	}

	if address == 0xFF44 {
		// todo: this is hardcoded for the doctor, but it shouldn't be
		return 0x90
	}

	fmt.Printf("Reading from %2.2X not supported (IO_REGISTERS)\n", address)
	return 0
}

func appendSerialToFile(value byte) {
	f, err := os.OpenFile("serial.out", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	bs, err := hex.DecodeString(fmt.Sprintf("%2.2X", value))
	if err != nil {
		panic(err)
	}

	if _, err = f.WriteString(string(bs)); err != nil {
		panic(err)
	}
}
