package goboy

import "fmt"

// @see https://gbdev.io/pandocs/Memory_Map.html
// 0x0000 - 0x3FFF : ROM Bank 00
// 0x4000 - 0x7FFF : ROM Bank 01 - switchable
// 0x8000 - 0x97FF : CHR RAM
// 0x9800 - 0x9BFF : BG Map 1
// 0x9C00 - 0x9FFF : BG Map 2
// 0xA000 - 0xBFFF : Cartridge RAM
// 0xC000 - 0xCFFF : RAM Bank 0
// 0xD000 - 0xDFFF : RAM Bank 1-7 - switchable - CGB only
// 0xE000 - 0xFDFF : Echo RAM (mirror of 0xC000 - 0xCFFF)
// 0xFE00 - 0xFE9F : Object Attribute Memory (OAM)
// 0xFEA0 - 0xFEFF : Not Usable
// 0xFF00 - 0xFF7F : I/O Registers
// 0xFF80 - 0xFFFE : High RAM (HRAM)
// 0xFFFF - 0xFFFF : Interrupt Enable Register

type Bus struct {
	cartridge Cartridge
}

func NewBus(cartridgePath string) Bus {
	return Bus{
		cartridge: LoadCartridge(cartridgePath),
	}
}

func (bus *Bus) read(address uint16) uint8 {
	if address < 0x4000 {
		return bus.cartridge.read(address)
	}

	panic(fmt.Sprintf("Reading from %2.2x not supported", address))
}

func (bus *Bus) write(address uint16, value byte) {
	if address < 0x4000 {
		bus.cartridge.write(address, value)
	}

	panic(fmt.Sprintf("Writing to %2.2x not supported", address))
}
