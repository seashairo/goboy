package goboy

import (
	"fmt"
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
	bus  *Bus
	vram *RAM
	oam  *RAM

	currentFrame      uint32
	scanlineTicks     uint32
	videoBuffer       [FRAME_BUFFER_SIZE]uint32
	previousFrameTime int64
	startTime         int64
	frameCount        int64
}

func NewPPU(bus *Bus) *PPU {
	return &PPU{
		bus:  bus,
		vram: NewRAM(8192, VIDEO_RAM_START),
		oam:  NewRAM(160, OAM_START),

		currentFrame:      0,
		scanlineTicks:     0,
		videoBuffer:       [FRAME_BUFFER_SIZE]uint32{},
		previousFrameTime: time.Now().UnixMilli(),
		startTime:         time.Now().UnixMilli(),
		frameCount:        0,
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

func (ppu *PPU) handleModeOam() {
	if ppu.scanlineTicks >= 80 {
		ppu.bus.io.lcd.SetMode(LCD_MODE_TRANSFER)
	}
}

func (ppu *PPU) handleModeTransfer() {
	if ppu.scanlineTicks >= 80+172 {
		ppu.bus.io.lcd.SetMode(LCD_MODE_HBLANK)
	}
}

func (ppu *PPU) handleModeVblank() {
	if ppu.scanlineTicks >= DOTS_PER_LINE {
		ppu.bus.io.lcd.IncrementLy()

		if ppu.bus.readByte(LY_ADDRESS) > SCANLINES_PER_FRAME {
			ppu.bus.io.lcd.SetMode(LCD_MODE_OAM)
			ppu.bus.writeByte(LY_ADDRESS, 0)
		}

		ppu.scanlineTicks = 0
	}
}

func (ppu *PPU) handleModeHblank() {
	if ppu.scanlineTicks >= DOTS_PER_LINE {
		ppu.bus.io.lcd.IncrementLy()

		if ppu.bus.readByte(LY_ADDRESS) >= LCD_HEIGHT {
			// If we're past the end of the screen, it's vblank time
			ppu.bus.io.lcd.SetMode(LCD_MODE_VBLANK)
			// The CPU has a specific vblank interrupt
			ppu.bus.interruptEnableRegister.SetInterrupt(INT_VBLANK, true)
			// And it the LCD wants, that can also trigger a stat interrupt
			if ppu.bus.io.lcd.CheckLcdStatusFlag(STAT_VBLANK_INTERRUPT) {
				ppu.bus.interruptEnableRegister.SetInterrupt(INT_LCD, true)
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
