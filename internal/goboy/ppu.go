package goboy

import "fmt"

type PPU struct {
	bus  *Bus
	vram *RAM
	oam  *RAM
}

func NewPPU(bus *Bus) *PPU {
	return &PPU{
		bus:  bus,
		vram: NewRAM(8192, VIDEO_RAM_START),
		oam:  NewRAM(160, OAM_START),
	}
}

func (ppu *PPU) readByte(address uint16) byte {
	if Between(address, VIDEO_RAM_START, VIDEO_RAM_END) {
		return ppu.vram.readByte(address)
	} else if Between(address, OAM_START, OAM_END) {
		return ppu.oam.readByte(address)
	}

	panic("Somehow didn't manage to read a byte")
}

func (ppu *PPU) writeByte(address uint16, value byte) {
	if Between(address, VIDEO_RAM_START, VIDEO_RAM_END) {
		ppu.vram.writeByte(address, value)
		return
	} else if Between(address, OAM_START, OAM_END) {
		ppu.oam.writeByte(address, value)
		return
	}

	panic(fmt.Sprintf("Somehow didn't manage to write a byte (%4.4X:%2.2X)\n", address, value))
}

func (ppu *PPU) Tick() {
}
