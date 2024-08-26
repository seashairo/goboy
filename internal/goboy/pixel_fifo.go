package goboy

type FetchState byte

const (
	FETCH_STATE_TILE = iota
	FETCH_STATE_DATA_LO
	FETCH_STATE_DATA_HI
	FETCH_STATE_SLEEP
	FETCH_STATE_PUSH
)

// @see https://gbdev.io/pandocs/pixel_fifo.html
type PixelFifo struct {
	bus *Bus

	data []uint32

	fetchState   FetchState
	lineX        byte
	pushedX      byte
	fetchX       byte
	bgwFetchData [3]byte
	oamData      [6]byte
	mapX         byte
	mapY         byte
	tileY        byte
	fifoX        byte
}

func NewPixelFifo(bus *Bus) *PixelFifo {
	return &PixelFifo{
		bus:          bus,
		data:         make([]uint32, 0),
		fetchState:   FETCH_STATE_TILE,
		lineX:        0,
		pushedX:      0,
		fetchX:       0,
		bgwFetchData: [3]byte{0, 0, 0},
		oamData:      [6]byte{0, 0, 0, 0, 0, 0},
		mapX:         0,
		mapY:         0,
		tileY:        0,
		fifoX:        0,
	}
}

func (pf *PixelFifo) push(pixel uint32) {
	pf.data = append(pf.data, pixel)
}

func (pf *PixelFifo) pop() uint32 {
	pixel := pf.data[0]
	pf.data = pf.data[1:]
	return pixel
}

func (pf *PixelFifo) SetState(state FetchState) {
	if state == FETCH_STATE_TILE {
		pf.fetchState = state
		pf.lineX = 0
		pf.fetchX = 0
		pf.pushedX = 0
		pf.fifoX = 0
	}
}

func (pf *PixelFifo) Process() {
	pf.mapX = (pf.fetchX + pf.bus.readByte(SCX_ADDRESS)) / 8
	pf.mapY = (pf.bus.readByte(LY_ADDRESS) + pf.bus.readByte(SCY_ADDRESS)) / 8
	pf.tileY = ((pf.bus.readByte(LY_ADDRESS) + pf.bus.readByte(SCY_ADDRESS)) % 8) * 2

	if pf.bus.ppu.scanlineTicks%2 == 0 {
		pf.Fetch()
	}

	pf.Push()
}

func (pf *PixelFifo) Push() {
	if len(pf.data) > 8 {
		data := pf.pop()

		if pf.lineX >= pf.bus.readByte(SCX_ADDRESS)%8 {
			index := uint32(pf.pushedX) + (uint32(pf.bus.readByte(LY_ADDRESS)) * LCD_WIDTH)
			pf.bus.ppu.videoBuffer[index] = data
			pf.pushedX += 1
		}

		pf.lineX += 1
	}
}

func (pf *PixelFifo) Fetch() {
	lcd := pf.bus.io.lcd

	switch pf.fetchState {
	case FETCH_STATE_TILE:
		if lcd.IsBgwEnabled() {
			pf.bgwFetchData[0] = pf.bus.readByte(lcd.BgTileMapOffset() + uint16(pf.mapX) + (uint16(pf.mapY) * 32))

			if lcd.BgwTileDataOffset() == 0x8800 {
				pf.bgwFetchData[0] += 128
			}
		}

		pf.fetchX += 8
		pf.fetchState = FETCH_STATE_DATA_LO
	case FETCH_STATE_DATA_LO:
		address := lcd.BgwTileDataOffset() + uint16(pf.bgwFetchData[0])*16 + uint16(pf.tileY)
		pf.bgwFetchData[1] = pf.bus.readByte(address)
		pf.fetchState = FETCH_STATE_DATA_HI
	case FETCH_STATE_DATA_HI:
		address := lcd.BgwTileDataOffset() + uint16(pf.bgwFetchData[0])*16 + uint16(pf.tileY) + 1
		pf.bgwFetchData[2] = pf.bus.readByte(address)
		pf.fetchState = FETCH_STATE_SLEEP
	case FETCH_STATE_SLEEP:
		pf.fetchState = FETCH_STATE_PUSH
	case FETCH_STATE_PUSH:
		if pf.add() {
			pf.fetchState = FETCH_STATE_TILE
		}
	}
}

func (pf *PixelFifo) add() bool {
	if len(pf.data) > 8 {
		return false
	}

	x := pf.fetchX - (pf.bus.readByte(SCX_ADDRESS) % 8)
	for i := 0; i < 8; i++ {
		bit := 7 - i

		hi := (pf.bgwFetchData[1] & (1 << bit)) >> bit
		lo := ((pf.bgwFetchData[2] & (1 << bit)) >> bit) << 1

		color := pf.bus.io.lcd.bgColors[hi|lo]

		if x > 0 {
			pf.push(color)
			pf.fifoX += 1
		}
	}

	return true
}

func (pf *PixelFifo) Reset() {
	pf.data = make([]uint32, 0)
}
