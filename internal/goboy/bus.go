package goboy

// @see https://gbdev.io/pandocs/Memory_Map.html
const (
	// 0x0000 - 0x3FFF : ROM Bank 00
	ROM_BANK_0_START = 0x0000
	ROM_BANK_0_END   = 0x3FFF

	// 0x4000 - 0x7FFF : ROM Bank 01 - switchable
	SWITCHABLE_ROM_BANK_START = 0x4000
	SWITCHABLE_ROM_BANK_END   = 0x7FFF

	// 0x8000 - 0x97FF : CHR RAM
	// 0x9800 - 0x9BFF : BG Map 1
	// 0x9C00 - 0x9FFF : BG Map 2
	VIDEO_RAM_START = 0x8000
	VIDEO_RAM_END   = 0x9FFF

	// 0xA000 - 0xBFFF : Cartridge RAM
	EXTERNAL_RAM_START = 0xA000
	EXTERNAL_RAM_END   = 0xBFFF

	// 0xC000 - 0xCFFF : RAM Bank 0
	WORK_RAM_START = 0xC000
	WORK_RAM_END   = 0xCFFF

	// 0xD000 - 0xDFFF : RAM Bank 1-7 - switchable - CGB only
	SWITCHABLE_WORK_RAM_START = 0xD000
	SWITCHABLE_WORK_RAM_END   = 0xDFFF

	// 0xE000 - 0xFDFF : Echo RAM (mirror of 0xC000 - 0xCFFF)
	ECHO_RAM_START = 0xE000
	ECHO_RAM_END   = 0xFDFF

	// 0xFE00 - 0xFE9F : Object Attribute Memory (OAM)
	OBJECT_ATTRIBUTE_MEMORY_START = 0xFE00
	OBJECT_ATTRIBUTE_MEMORY_END   = 0xFE9F

	// 0xFEA0 - 0xFEFF : Not Usable
	NOT_USABLE_START = 0xFEA0
	NOT_USABLE_END   = 0xFEFF

	// 0xFF00 - 0xFF7F : I/O Registers
	IO_REGISTERS_START = 0xFF00
	IO_REGISTERS_END   = 0xFF7F

	// 0xFF80 - 0xFFFE : High RAM (HRAM)
	HIGH_RAM_START = 0xFF80
	HIGH_RAM_END   = 0xFFFE

	// 0xFFFF - 0xFFFF : Interrupt Enable Register
	INTERRUPT_ENABLE_REGISTER_START = 0xFFFF
	INTERRUPT_ENABLE_REGISTER_END   = 0xFFFF
)

type Bus struct {
	cartridge               Cartridge
	interruptEnableRegister InterruptRegister
	wram                    RAM
	hram                    RAM
	io                      IO
}

func NewBus(cartridgePath string) Bus {
	return Bus{
		cartridge:               LoadCartridge(cartridgePath),
		interruptEnableRegister: NewInterruptRegister(0),
		wram:                    NewRAM(SWITCHABLE_WORK_RAM_END-WORK_RAM_START+1, WORK_RAM_START),
		hram:                    NewRAM(HIGH_RAM_END-HIGH_RAM_START+1, HIGH_RAM_START),
		io:                      NewIO(),
	}
}

func (bus *Bus) readByte(address uint16) byte {
	if address <= SWITCHABLE_ROM_BANK_END {
		return bus.cartridge.readByte(address)
	} else if address <= VIDEO_RAM_END {
		// fmt.Printf("Reading from %2.2X not supported (VIDEO_RAM)\n", address)
	} else if address <= EXTERNAL_RAM_END {
		return bus.cartridge.readByte(address)
	} else if address <= SWITCHABLE_WORK_RAM_END {
		return bus.wram.readByte(address)
	} else if address <= ECHO_RAM_END {
		// fmt.Printf("Reading from %2.2X not supported (ECHO_RAM)\n", address)
	} else if address <= OBJECT_ATTRIBUTE_MEMORY_END {
		// fmt.Printf("Reading from %2.2X not supported (OBJECT_ATTRIBUTE_MEMORY)\n", address)
	} else if address <= NOT_USABLE_END {
		// fmt.Printf("Reading from %2.2X not supported (NOT_USABLE)\n", address)
	} else if address <= IO_REGISTERS_END {
		return bus.io.readByte(address)
	} else if address <= HIGH_RAM_END {
		return bus.hram.readByte(address)
	}

	return bus.interruptEnableRegister.readByte()
}

func (bus *Bus) readWord(address uint16) uint16 {
	lo := bus.readByte(address)
	hi := bus.readByte(address + 1)

	return BytesToUint16(hi, lo)
}

func (bus *Bus) writeByte(address uint16, value byte) {
	if address <= SWITCHABLE_ROM_BANK_END {
		bus.cartridge.writeByte(address, value)
	} else if address <= VIDEO_RAM_END {
		// fmt.Printf("Writing to %2.2X not supported (VIDEO_RAM)\n", address)
	} else if address <= EXTERNAL_RAM_END {
		bus.cartridge.writeByte(address, value)
	} else if address <= SWITCHABLE_WORK_RAM_END {
		bus.wram.writeByte(address, value)
	} else if address <= ECHO_RAM_END {
		// fmt.Printf("Writing to %2.2X not supported (ECHO_RAM)\n", address)
	} else if address <= OBJECT_ATTRIBUTE_MEMORY_END {
		// fmt.Printf("Writing to %2.2X not supported (OBJECT_ATTRIBUTE_MEMORY)\n", address)
	} else if address <= NOT_USABLE_END {
		// fmt.Printf("Writing to %2.2X not supported (NOT_USABLE)\n", address)
	} else if address <= IO_REGISTERS_END {
		bus.io.writeByte(address, value)
	} else if address <= HIGH_RAM_END {
		bus.hram.writeByte(address, value)
	}

	bus.interruptEnableRegister.writeByte(value)
}

func (bus *Bus) writeWord(address uint16, value uint16) {
	hi, lo := Uint16ToBytes(value)
	bus.writeByte(address, lo)
	bus.writeByte(address+1, hi)
}
