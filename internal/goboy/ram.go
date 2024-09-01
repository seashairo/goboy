package goboy

import "fmt"

type RAM struct {
	// The size of the RAM bank in bytes
	size uint32

	// The starting address of the RAM bank (e.g. 0xA000 for WRAM)
	offset uint16

	// The raw data in the RAM bank
	data []byte
}

func NewRAM(size uint32, offset uint16) *RAM {
	if DEBUG {
		fmt.Printf("Creating RAM bank with size 0x%4.4X and offset 0x%4.4X\n", size, offset)
	}

	return &RAM{
		size:   size,
		offset: offset,
		data:   make([]byte, size),
	}
}

func (ram *RAM) readByte(address uint16) byte {
	return ram.data[address-ram.offset]
}

func (ram *RAM) writeByte(address uint16, value byte) {
	ram.data[address-ram.offset] = value
}

func (ram *RAM) debugPrint() {
	const bytesPerRow = 16

	for i := 0; i < len(ram.data); i += bytesPerRow {
		out := fmt.Sprintf("%04X: ", i+int(ram.offset))

		// Print the hex values
		for j := 0; j < bytesPerRow && i+j < len(ram.data); j++ {
			out += fmt.Sprintf("%02X ", ram.data[i+j])
		}

		// Print spacing between hex values and ASCII characters
		for j := len(ram.data[i:]); j < bytesPerRow; j++ {
			out += "   "
		}

		// Print ASCII characters (if printable)
		for j := 0; j < bytesPerRow && i+j < len(ram.data); j++ {
			b := ram.data[i+j]
			if b >= 32 && b <= 126 { // Printable ASCII range
				out += fmt.Sprintf("%c", b)
			} else {
				out += "."
			}
		}

		fmt.Println(out)
	}
}
