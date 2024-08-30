package goboy

import (
	"fmt"
	"slices"
	"time"
)

// @see https://gbdev.io/pandocs/Rendering.html
const (
	SCANLINES_PER_FRAME = 154
	DOTS_PER_LINE       = 456
	FRAME_BUFFER_SIZE   = LCD_HEIGHT * LCD_WIDTH
	LIMIT_FPS           = true
	TARGET_FRAME_TIME   = int64(1000 / 60) // 60 fps
)

type OamEntryFlag byte

const (
	OAM_VRAM_BANK = iota + 3
	OAM_DMG_PALETTE
	OAM_X_FLIP
	OAM_Y_FLIP
	OAM_PRIORITY
)

// @see https://gbdev.io/pandocs/OAM.html
type OamEntry struct {
	y     byte
	x     byte
	tile  byte
	flags byte
}

func (o OamEntry) Check(flag OamEntryFlag) bool {
	return GetBit(o.flags, byte(flag))
}

type PPU struct {
	bus       *Bus
	vram      *RAM
	oam       *[40]OamEntry
	pixelFifo *PixelFifo

	lineSprites       []OamEntry
	windowLine        uint32
	currentFrame      uint32
	scanlineTicks     uint32
	videoBuffer       *[FRAME_BUFFER_SIZE]uint32
	previousFrameTime int64
	startTime         int64
	frameCount        int64
}

func NewPPU(bus *Bus) *PPU {
	return &PPU{
		bus:       bus,
		vram:      NewRAM(8192, VIDEO_RAM_START),
		oam:       &[40]OamEntry{},
		pixelFifo: NewPixelFifo(bus),

		// sprites
		lineSprites: make([]OamEntry, 40),
		windowLine:  0,

		// rendering
		currentFrame:  0,
		scanlineTicks: 0,
		videoBuffer:   &[FRAME_BUFFER_SIZE]uint32{},

		// fps
		previousFrameTime: time.Now().UnixMilli(),
		startTime:         time.Now().UnixMilli(),
		frameCount:        0,
	}
}

func (ppu *PPU) readByte(address uint16) byte {
	if Between(address, VIDEO_RAM_START, VIDEO_RAM_END) {
		return ppu.vram.readByte(address)
	} else if Between(address, OAM_START, OAM_END) {
		oamEntry := ppu.getOamEntry(address)

		switch address % 4 {
		case 0:
			return oamEntry.y
		case 1:
			return oamEntry.x
		case 2:
			return oamEntry.tile
		case 3:
			return oamEntry.flags
		}
	}

	panic("Somehow didn't manage to read a byte")
}

func (ppu *PPU) writeByte(address uint16, value byte) {
	if Between(address, VIDEO_RAM_START, VIDEO_RAM_END) {
		ppu.vram.writeByte(address, value)
		return
	} else if Between(address, OAM_START, OAM_END) {
		oamEntry := ppu.getOamEntry(address)

		switch address % 4 {
		case 0:
			oamEntry.y = value
		case 1:
			oamEntry.x = value
		case 2:
			oamEntry.tile = value
		case 3:
			oamEntry.flags = value
		}

		return
	}

	panic(fmt.Sprintf("Somehow didn't manage to write a byte (%4.4X:%2.2X)\n", address, value))
}

func (ppu *PPU) getOamEntry(address uint16) *OamEntry {
	offset := (address - OAM_START)
	oamIndex := offset / 4
	return &ppu.oam[oamIndex]
}

func (ppu *PPU) Tick() {
	ppu.scanlineTicks++

	switch ppu.bus.io.lcd.GetMode() {
	case LCD_MODE_HBLANK:
		ppu.handleModeHblank()
	case LCD_MODE_VBLANK:
		ppu.handleModeVblank()
	case LCD_MODE_OAM:
		ppu.handleModeOam()
	case LCD_MODE_TRANSFER:
		ppu.handleModeTransfer()
	}
}

func (ppu *PPU) loadLineSprites() {
	ppu.lineSprites = nil

	// This is the line we're fetching sprites for
	ly := ppu.bus.readByte(LY_ADDRESS)
	spriteHeight := ppu.bus.io.lcd.ObjSize()

	for i := 0; i < len(ppu.oam); i++ {
		sprite := ppu.oam[i]

		// A sprite is on the line if it starts on or above the scanline, and
		// finishes below the scanline. The sprite Y coordinates are always offset
		// by 16 (I don't know why - see https://gbdev.io/pandocs/OAM.html)
		if sprite.y-16 <= ly && sprite.y-16+spriteHeight > ly {
			ppu.lineSprites = append(ppu.lineSprites, sprite)
		}
	}

	slices.SortFunc(ppu.lineSprites, func(a, b OamEntry) int {
		if a.x < b.x {
			return -1
		}

		if a.x > b.x {
			return 1
		}

		return 0
	})

	// The Game Boy will only render 10 sprites per line, so if we have more than
	// 10 sprites in the list, we sort them by X position and only keep the first
	// 10. There is a tie breaker for index, so I might need to revisit this to
	// make sure it's got the right sprites.
	if len(ppu.lineSprites) > 10 {
		ppu.lineSprites = ppu.lineSprites[:10]
	}
}

func (ppu *PPU) handleModeOam() {
	// I don't know when sprite data is actually loaded, but we probably only need
	// to do it once during OAM phase and not 80 times
	if ppu.scanlineTicks == 1 {
		ppu.loadLineSprites()
	}

	// After 80 ticks on this line, we move to mode 3 and start pushing data
	if ppu.scanlineTicks >= 80 {
		ppu.bus.io.lcd.SetMode(LCD_MODE_TRANSFER)

		ppu.pixelFifo.fetchState = FETCH_STATE_TILE
		ppu.pixelFifo.lineX = 0
		ppu.pixelFifo.fetchX = 0
		ppu.pixelFifo.pushedX = 0
		ppu.pixelFifo.fifoX = 0
	}
}

func (ppu *PPU) handleModeTransfer() {
	ppu.pixelFifo.Process()

	if ppu.pixelFifo.pushedX >= LCD_WIDTH {
		ppu.pixelFifo.Reset()
		ppu.bus.io.lcd.SetMode(LCD_MODE_HBLANK)

		if ppu.bus.io.lcd.CheckLcdStatusFlag(STAT_HBLANK_INTERRUPT) {
			ppu.bus.io.interrupts.SetInterrupt(INT_LCD, true)
		}
	}
}

func (ppu *PPU) incrementLy() {
	ly := ppu.bus.readByte(LY_ADDRESS)
	wy := ppu.bus.readByte(WY_ADDRESS)

	if ppu.isWindowVisible() && ly >= wy && ly < wy+LCD_HEIGHT {
		ppu.windowLine += 1
	}

	ppu.bus.io.lcd.IncrementLy()
}

func (ppu *PPU) isWindowVisible() bool {
	wx := ppu.bus.readByte(WX_ADDRESS)
	wy := ppu.bus.readByte(WY_ADDRESS)

	return ppu.bus.io.lcd.IsWindowEnabled() && wx <= 166 && wy < LCD_HEIGHT
}

func (ppu *PPU) handleModeVblank() {
	if ppu.scanlineTicks >= DOTS_PER_LINE {
		ppu.incrementLy()

		if ppu.bus.readByte(LY_ADDRESS) >= SCANLINES_PER_FRAME {
			ppu.bus.io.lcd.SetMode(LCD_MODE_OAM)
			ppu.bus.writeByte(LY_ADDRESS, 0)
			ppu.windowLine = 0
		}

		ppu.scanlineTicks = 0
	}
}

func (ppu *PPU) handleModeHblank() {
	if ppu.scanlineTicks >= DOTS_PER_LINE {
		ppu.incrementLy()

		if ppu.bus.readByte(LY_ADDRESS) >= LCD_HEIGHT {
			// If we're past the end of the screen, it's vblank time
			ppu.bus.io.lcd.SetMode(LCD_MODE_VBLANK)
			// The CPU has a specific vblank interrupt
			ppu.bus.io.interrupts.SetInterrupt(INT_VBLANK, true)
			// And if the LCD wants, that can also trigger a stat interrupt
			if ppu.bus.io.lcd.CheckLcdStatusFlag(STAT_VBLANK_INTERRUPT) {
				ppu.bus.io.interrupts.SetInterrupt(INT_LCD, true)
			}
			ppu.currentFrame++

			if LIMIT_FPS {
				currentTime := time.Now().UnixMilli()
				frameTime := currentTime - ppu.previousFrameTime

				if frameTime < TARGET_FRAME_TIME {
					// todo: this isn't quite right, it's hitting 70 fps, try moving timer
					// to its own class
					time.Sleep(time.Duration(TARGET_FRAME_TIME-frameTime) * time.Millisecond)
				}

				if currentTime-ppu.startTime > 1000 {
					fps := ppu.frameCount
					ppu.startTime = currentTime
					ppu.frameCount = 0
					fmt.Printf("fps: %d\n", fps)
				}

				ppu.frameCount++
				ppu.previousFrameTime = currentTime
			}
		} else {
			ppu.bus.io.lcd.SetMode(LCD_MODE_OAM)
		}

		ppu.scanlineTicks = 0
	}
}
