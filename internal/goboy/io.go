package goboy

import (
	"encoding/hex"
	"fmt"
	"os"
)

type IO struct {
	interrupts *InterruptRegister
	timer      *Timer
	dma        *DMA
}

func NewIO(bus *Bus, timer *Timer, interruptEnableRegister *InterruptRegister) *IO {
	// todo: should all these things be singletons? the architecture needs to
	// change so different "devices" all have access to the bus and are also
	// accessible through the bus
	dma := NewDMA(bus)

	return &IO{
		interrupts: interruptEnableRegister,
		timer:      timer,
		dma:        dma,
	}
}

func (io *IO) writeByte(address uint16, value byte) {
	if address == 0xFF01 {
		appendSerialToFile(value)
		return
	}

	if address >= 0xFF04 && address <= 0xFF07 {
		io.timer.writeByte(address, value)
		return
	}

	if address == 0xFF0F {
		io.interrupts.writeByte(value)
		return
	}

	if address == 0xFF46 {
		io.dma.writeByte(address, value)
	}

	// fmt.Printf("Writing to %2.2X not supported (IO_REGISTERS)\n", address)
}

var ly = 0

func (io *IO) readByte(address uint16) byte {
	if address >= 0xFF04 && address <= 0xFF07 {
		io.timer.readByte(address)
	}

	if address == 0xFF0F {
		return io.interrupts.readByte()
	}

	if address == 0xFF44 {
		// todo: this is hardcoded for the doctor, but it shouldn't be
		// return 0x90

		ly++
		return byte(ly)
	}

	if address == 0xFF46 {
		return io.dma.readByte(address)
	}

	// fmt.Printf("Reading from %2.2X not supported (IO_REGISTERS)\n", address)
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
