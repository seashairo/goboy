package goboy

import "fmt"

type RAM struct {
	// The size of the RAM bank in bytes
	size uint16

	// The starting address of the RAM bank (e.g. 0xA000 for WRAM)
	offset uint16

	// The raw data in the RAM bank
	data []byte
}

func NewRAM(size uint16, offset uint16) RAM {
	fmt.Printf("Creating RAM bank with size 0x%4.4X and offset 0x%4.4X\n", size, offset)

	return RAM{
		size:   size,
		offset: offset,
		data:   make([]byte, size),
	}
}

func (ram *RAM) readByte(address uint16) byte {
	return ram.data[address-ram.offset]
}

func (ram *RAM) readWord(address uint16) uint16 {
	lo := ram.readByte(address)
	hi := ram.readByte(address + 1)

	return BytesToUint16(hi, lo)
}

func (ram *RAM) writeByte(address uint16, value byte) {
	ram.data[address-ram.offset] = value
}

func (ram *RAM) writeWord(address uint16, value uint16) {
	hi, lo := Uint16ToBytes(value)
	ram.writeByte(address, lo)
	ram.writeByte(address+1, hi)
}
