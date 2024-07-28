package goboy

type PPU struct {
	bus *Bus
}

func NewPPU(bus *Bus) PPU {
	return PPU{
		bus: bus,
	}
}

func (ppu *PPU) Tick() {
}
