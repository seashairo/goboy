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
	ppu *PPU
	lcd *LCD

	data []uint32

	fetchState        FetchState
	lineX             byte
	pushedX           byte
	fetchX            byte
	bgwFetchData      [3]byte
	oamFetchData      [6]byte
	fetchedOamEntries []OamEntry
	mapX              byte
	mapY              byte
	tileY             byte
	fifoX             byte
}

func NewPixelFifo(bus *Bus, ppu *PPU, lcd *LCD) *PixelFifo {
	return &PixelFifo{
		bus:               bus,
		ppu:               ppu,
		lcd:               lcd,
		data:              make([]uint32, 0),
		fetchState:        FETCH_STATE_TILE,
		lineX:             0,
		pushedX:           0,
		fetchX:            0,
		bgwFetchData:      [3]byte{0, 0, 0},
		oamFetchData:      [6]byte{0, 0, 0, 0, 0, 0},
		fetchedOamEntries: make([]OamEntry, 0),
		mapX:              0,
		mapY:              0,
		tileY:             0,
		fifoX:             0,
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
	pf.mapX = (pf.fetchX + pf.bus.readByte(LCD_SCX)) / 8
	pf.mapY = (pf.bus.readByte(LCD_LY) + pf.bus.readByte(LCD_SCY)) / 8
	pf.tileY = ((pf.bus.readByte(LCD_LY) + pf.bus.readByte(LCD_SCY)) % 8) * 2

	if pf.ppu.scanlineTicks%2 == 0 {
		pf.Fetch()
	}

	pf.Push()
}

func (pf *PixelFifo) Push() {
	if len(pf.data) > 8 {
		data := pf.pop()

		if pf.lineX >= pf.bus.readByte(LCD_SCX)%8 {
			index := uint32(pf.pushedX) + (uint32(pf.bus.readByte(LCD_LY)) * LCD_WIDTH)
			pf.ppu.videoBuffer[index] = data
			pf.pushedX += 1
		}

		pf.lineX += 1
	}
}

func (pf *PixelFifo) Fetch() {
	lcd := pf.lcd

	switch pf.fetchState {
	case FETCH_STATE_TILE:
		pf.fetchedOamEntries = nil

		if lcd.IsBgwEnabled() {
			pf.bgwFetchData[0] = pf.bus.readByte(lcd.BgTileMapOffset() + uint16(pf.mapX) + (uint16(pf.mapY) * 32))

			if lcd.BgwTileDataOffset() == 0x8800 {
				pf.bgwFetchData[0] += 128
			}

			pf.loadWindowTile()
		}

		if lcd.IsObjEnabled() && len(pf.ppu.lineSprites) != 0 {
			pf.loadSpriteTile()
		}

		pf.fetchX += 8
		pf.fetchState = FETCH_STATE_DATA_LO
	case FETCH_STATE_DATA_LO:
		address := lcd.BgwTileDataOffset() + uint16(pf.bgwFetchData[0])*16 + uint16(pf.tileY)
		pf.bgwFetchData[1] = pf.bus.readByte(address)
		pf.loadSpriteData(0)
		pf.fetchState = FETCH_STATE_DATA_HI
	case FETCH_STATE_DATA_HI:
		address := lcd.BgwTileDataOffset() + uint16(pf.bgwFetchData[0])*16 + uint16(pf.tileY) + 1
		pf.bgwFetchData[2] = pf.bus.readByte(address)
		pf.loadSpriteData(1)
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

	x := pf.fetchX - (pf.bus.readByte(LCD_SCX) % 8)
	for i := 0; i < 8; i++ {
		bit := 7 - i

		lo := (pf.bgwFetchData[1] & (1 << bit)) >> bit
		hi := ((pf.bgwFetchData[2] & (1 << bit)) >> bit) << 1
		colorIndex := hi | lo

		color := pf.lcd.bgColors[colorIndex]

		if !pf.lcd.IsBgwEnabled() {
			color = pf.lcd.bgColors[0]
		}

		if pf.lcd.IsObjEnabled() {
			color = pf.fetchSpritePixels(color, colorIndex)
		}

		if x > 0 {
			pf.push(color)
			pf.fifoX += 1
		}
	}

	return true
}

func (pf *PixelFifo) loadSpriteData(offset int) {
	ly := pf.bus.readByte(LCD_LY)
	spriteHeight := pf.lcd.ObjSize()

	for i := 0; i < len(pf.fetchedOamEntries); i++ {
		sprite := pf.fetchedOamEntries[i]
		tileY := (ly + 16 - sprite.y) * 2

		if sprite.Check(OAM_Y_FLIP) {
			tileY = (spriteHeight*2 - 2) - tileY
		}

		tileIndex := sprite.tile
		if spriteHeight == 16 {
			tileIndex &= 0b11111110
		}

		address := VIDEO_RAM_START + uint16(tileIndex)*16 + uint16(tileY) + uint16(offset)
		pf.oamFetchData[(i*2)+offset] = pf.bus.readByte(address)
	}
}

func (pf *PixelFifo) loadSpriteTile() {
	for i := 0; i < len(pf.ppu.lineSprites); i++ {
		sprite := pf.ppu.lineSprites[i]
		spriteX := sprite.x - 8 + pf.bus.readByte(LCD_SCX)%8

		if (spriteX >= pf.fetchX && spriteX < pf.fetchX+8) ||
			((spriteX+8) >= pf.fetchX && (spriteX) < pf.fetchX) {
			pf.fetchedOamEntries = append(pf.fetchedOamEntries, sprite)
		}

		if len(pf.fetchedOamEntries) >= 3 {
			break
		}
	}
}

func (pf *PixelFifo) loadWindowTile() {
	if !pf.ppu.isWindowVisible() {
		return
	}

	wx := pf.bus.readByte(LCD_WX)
	wy := pf.bus.readByte(LCD_WY)
	ly := pf.bus.readByte(LCD_LY)

	fetchX := pf.fetchX + 7

	if fetchX >= wx &&
		fetchX < wx+LCD_WIDTH {
		if ly >= wy && ly < wy+LCD_HEIGHT {
			tx := (fetchX - wx) / 8
			ty := pf.ppu.windowLine / 8

			base := pf.lcd.WindowTileMapOffset()
			address := base + uint16(tx) + uint16(ty)*32

			pf.bgwFetchData[0] = pf.bus.readByte(address)
			if pf.lcd.BgwTileDataOffset() == 0x8800 {
				pf.bgwFetchData[0] += 128
			}
		}
	}
}

func (pf *PixelFifo) fetchSpritePixels(bgColor uint32, bgColorIndex byte) uint32 {
	color := bgColor

	for i := 0; i < len(pf.fetchedOamEntries); i++ {
		sprite := pf.fetchedOamEntries[i]
		spriteX := sprite.x - 8 + pf.bus.readByte(LCD_SCX)%8

		if spriteX+8 < pf.fifoX {
			continue
		}

		offsetX := pf.fifoX - spriteX
		if offsetX > 7 {
			continue
		}

		bit := 7 - int(offsetX)
		if sprite.Check(OAM_X_FLIP) {
			bit = int(offsetX)
		}

		lo := (pf.oamFetchData[i*2] & (1 << bit)) >> bit
		hi := ((pf.oamFetchData[i*2+1] & (1 << bit)) >> bit) << 1
		colorIndex := hi | lo

		// For sprites, if they would use the first entry in the palette they are
		// instead considered to be transparent so we'll let the buffer keep the BGW
		// color. The same happens if the sprite gives priority to the background
		// and the background isn't a blank space
		if colorIndex == 0 || (sprite.Check(OAM_PRIORITY) && bgColorIndex != 0) {
			continue
		}

		// The DMG has two different palettes that could be in use depending on this
		// flag so we need to make sure we pull the right one
		color = pf.lcd.sp1Colors[colorIndex]
		if sprite.Check(OAM_DMG_PALETTE) {
			color = pf.lcd.sp2Colors[colorIndex]
		}

		// Sprites are only mixed if they're transparent, so if we've gotten this
		// far we take the color and break
		break
	}

	return color
}

func (pf *PixelFifo) Reset() {
	pf.data = nil
}
