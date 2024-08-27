package goboy

import "fmt"

// @see https://gbdev.io/pandocs/OAM_DMA_Transfer.html#ff46--dma-oam-dma-source-address--start
type DMA struct {
	bus *Bus

	// Whether transfer is in progress. An active DMA transfer prevents
	// certain operations from operating
	active bool
	// DMA takes 2 cycles to start, so there might be a couple of ticks to wait
	// before transferring
	delay byte
	// The hi byte of the address being written to
	addressHi uint16
	// The index of the byte being copied
	byteIndex uint16
}

func NewDMA(bus *Bus) *DMA {
	return &DMA{
		bus:       bus,
		active:    false,
		delay:     2,
		addressHi: 0x0000,
		byteIndex: 0,
	}
}

func (dma *DMA) writeByte(address uint16, value byte) {
	dma.Start(value)
}

func (dma *DMA) readByte(address uint16) byte {
	return byte(dma.addressHi)
}

func (dma *DMA) Start(addressHi byte) {
	dma.active = true
	dma.delay = 2
	dma.addressHi = uint16(addressHi)
	dma.byteIndex = 0

	fmt.Printf("Starting DMA from %4.4X\n", dma.addressHi*0x100)
}

func (dma *DMA) Tick() {
	if !dma.active {
		return
	}

	if dma.delay > 0 {
		dma.delay--
		return
	}

	srcAddress := (dma.addressHi * 0x100) + dma.byteIndex
	dstAddress := uint16(OAM_START) + dma.byteIndex
	dma.bus.writeByte(dstAddress, dma.bus.readByte(srcAddress))

	dma.byteIndex++
	dma.active = dma.byteIndex < 0xA0
}

func (dma *DMA) Active() bool {
	return dma.active
}
