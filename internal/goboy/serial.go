package goboy

import (
	"encoding/hex"
	"fmt"
	"os"
)

const (
	SERIAL_SB = 0xFF01
	SERIAL_SC = 0xFF02
)

const (
	SC_CLOCK_SELECT    = 0
	SC_CLOCK_SPEED     = 1
	SC_TRANSFER_ENABLE = 7
)

// @see https://gbdev.io/pandocs/Serial_Data_Transfer_(Link_Cable).html
type Serial struct {
	gameboy *GameBoy

	sc byte
	sb byte

	transferredBits byte
	outgoingByte    byte
}

func NewSerial(gameboy *GameBoy) *Serial {
	return &Serial{
		gameboy:         gameboy,
		sc:              0,
		sb:              0,
		transferredBits: 0,
		outgoingByte:    0,
	}
}

func (serial *Serial) Tick() {
	if serial.useInternalClock() {
		// For the original GameBoy, the Serial clock is half the speed of the system
		// clock, so we can return early half the time
		if serial.gameboy.timer.sysclk%2 == 0 {
			return
		}
	} else {
		// If we're waiting for an external clock, we will never get one. There is
		// currently no actual serial support, so we'll wait forever for an external
		// clock to appear
		return
	}

	if !serial.transferEnabled() {
		return
	}

	outgoingBit := GetBit(serial.sb, 7)
	serial.outgoingByte = SetBit(serial.outgoingByte<<1, 0, outgoingBit)

	serial.sb = (serial.sb << 1) | 1

	serial.transferredBits++

	if serial.transferredBits != 8 {
		return
	}

	appendSerialToFile(serial.outgoingByte)
	serial.transferredBits = 0
	serial.outgoingByte = 0

	serial.sc = SetBit(serial.sc, SC_TRANSFER_ENABLE, false)
	serial.gameboy.RequestInterrupt(INT_SERIAL)
}

func (serial *Serial) transferEnabled() bool {
	return GetBit(serial.sc, SC_TRANSFER_ENABLE)
}

func (serial *Serial) useInternalClock() bool {
	return GetBit(serial.sc, SC_CLOCK_SELECT)
}

func (serial *Serial) readByte(address uint16) byte {
	switch address {
	case SERIAL_SB:
		return serial.sb
	case SERIAL_SC:
		return serial.sc
	}

	return 0
}

func (serial *Serial) writeByte(address uint16, value byte) {
	switch address {
	case SERIAL_SB:
		serial.sb = value
	case SERIAL_SC:
		serial.sc = value
	}
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
