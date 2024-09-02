package goboy

import "fmt"

const (
	LCD_WIDTH  = 160
	LCD_HEIGHT = 144
)

type LcdControlFlag byte

const (
	BGW_ENABLE = iota
	OBJ_ENABLE
	OBJ_SIZE
	BG_TILE_MAP
	BGW_TILES
	WINDOW_ENABLE
	WINDOW_TILE_MAP
	LCD_ENABLE
)

type LcdStatusFlag byte

const (
	STAT_LYC_EQUAL = iota + 2
	STAT_HBLANK_INTERRUPT
	STAT_VBLANK_INTERRUPT
	STAT_OAM_INTERRUPT
	STAT_LYC_INTERRUPT
)

type LcdMode byte

const (
	LCD_MODE_HBLANK = iota
	LCD_MODE_VBLANK
	LCD_MODE_OAM
	LCD_MODE_TRANSFER
)

const (
	LCD_LCDC uint16 = 0xFF40
	LCD_STAT uint16 = 0xFF41
	LCD_LY   uint16 = 0xFF44
	LCD_LYC  uint16 = 0xFF45
	LCD_SCY  uint16 = 0xFF42
	LCD_SCX  uint16 = 0xFF43
	LCD_WX   uint16 = 0xFF4B
	LCD_WY   uint16 = 0xFF4A
	LCD_BGP  uint16 = 0xFF47
	LCD_OBJ0 uint16 = 0xFF48
	LCD_OBJ1 uint16 = 0xFF49
)

var DEFAULT_PALLETTE = [4]uint32{0xFFFFFFFF, 0xFFA9A9A9, 0xFF545454, 0xFF000000}

type LCD struct {
	gameboy *GameBoy
	bus     *Bus

	// @see https://gbdev.io/pandocs/LCDC.html
	lcdc byte // 0xFF40
	// @see https://gbdev.io/pandocs/STAT.html
	stat byte // 0xFF41
	ly   byte // 0xFF44
	lyc  byte // 0xFF45
	// @see https://gbdev.io/pandocs/Scrolling.html#ff42ff43--scy-scx-background-viewport-y-position-x-position
	scy byte // 0xFF42
	scx byte // 0xFF43
	wx  byte // 0xFF4A
	wy  byte // 0xFF4B
	// @see https://gbdev.io/pandocs/Palettes.html
	bgp  byte // 0xFF47
	obj0 byte // 0xFF48
	obj1 byte // 0xFF49

	bgColors  [4]uint32
	sp1Colors [4]uint32
	sp2Colors [4]uint32
}

func NewLCD(gameboy *GameBoy, bus *Bus) *LCD {
	return &LCD{
		gameboy:   gameboy,
		bus:       bus,
		lcdc:      0x91,
		stat:      0,
		ly:        0,
		lyc:       0,
		scy:       0,
		scx:       0,
		wx:        0,
		wy:        0,
		bgp:       0xFC,
		obj0:      0xFF,
		obj1:      0xFF,
		bgColors:  DEFAULT_PALLETTE,
		sp1Colors: DEFAULT_PALLETTE,
		sp2Colors: DEFAULT_PALLETTE,
	}
}

func (lcd *LCD) readByte(address uint16) byte {
	switch address {
	case LCD_LCDC:
		return lcd.lcdc
	case LCD_STAT:
		return lcd.stat
	case LCD_LY:
		return lcd.ly
	case LCD_LYC:
		return lcd.lyc
	case LCD_SCY:
		return lcd.scy
	case LCD_SCX:
		return lcd.scx
	case LCD_WX:
		return lcd.wx
	case LCD_WY:
		return lcd.wy
	case LCD_BGP:
		return lcd.bgp
	case LCD_OBJ0:
		return lcd.obj0
	case LCD_OBJ1:
		return lcd.obj1
	}

	panic(fmt.Sprintf("Bad read from LCD: %4.4X", address))
}

func (lcd *LCD) writeByte(address uint16, value byte) {
	switch address {
	case LCD_LCDC:
		lcd.lcdc = value
		return
	case LCD_STAT:
		lcd.stat = value
		return
	case LCD_LY:
		lcd.ly = value
		return
	case LCD_LYC:
		lcd.lyc = value
		return
	case LCD_SCY:
		lcd.scy = value
		return
	case LCD_SCX:
		lcd.scx = value
		return
	case LCD_WX:
		lcd.wx = value
		return
	case LCD_WY:
		lcd.wy = value
		return
	case LCD_BGP:
		lcd.updatePalette(&lcd.bgColors, value)
		lcd.bgp = value
		return
	case LCD_OBJ0:
		lcd.updatePalette(&lcd.sp1Colors, value&0b11111100)
		lcd.obj0 = value
		return
	case LCD_OBJ1:
		lcd.updatePalette(&lcd.sp2Colors, value&0b11111100)
		lcd.obj1 = value
		return
	}

	panic(fmt.Sprintf("Bad write to LCD: %4.4X", address))
}

func (lcd *LCD) SetLcdControlFlag(flag LcdControlFlag, on bool) {
	lcd.lcdc = SetBit(lcd.lcdc, byte(flag), on)
}

func (lcd *LCD) CheckLcdControlFlag(flag LcdControlFlag) bool {
	return GetBit(lcd.lcdc, byte(flag))
}

func (lcd *LCD) IsBgwEnabled() bool {
	return GetBit(lcd.lcdc, byte(BGW_ENABLE))
}

func (lcd *LCD) IsObjEnabled() bool {
	return GetBit(lcd.lcdc, byte(OBJ_ENABLE))
}

func (lcd *LCD) ObjSize() byte {
	if GetBit(lcd.lcdc, byte(OBJ_SIZE)) {
		return 16
	}

	return 8
}

func (lcd *LCD) BgTileMapOffset() uint16 {
	if GetBit(lcd.lcdc, byte(BG_TILE_MAP)) {
		return 0x9C00
	}

	return 0x9800
}

func (lcd *LCD) BgwTileDataOffset() uint16 {
	if GetBit(lcd.lcdc, byte(BGW_TILES)) {
		return 0x8000
	}

	return 0x8800
}

func (lcd *LCD) IsWindowEnabled() bool {
	return GetBit(lcd.lcdc, byte(WINDOW_ENABLE))
}

func (lcd *LCD) WindowTileMapOffset() uint16 {
	if GetBit(lcd.lcdc, byte(WINDOW_TILE_MAP)) {
		return 0x9C00
	}

	return 0x9800
}

func (lcd *LCD) IsLcdEnabled() bool {
	return GetBit(lcd.lcdc, byte(LCD_ENABLE))
}

func (lcd *LCD) SetLcdStatusFlag(flag LcdStatusFlag, on bool) {
	lcd.stat = SetBit(lcd.stat, byte(flag), on)
}

func (lcd *LCD) CheckLcdStatusFlag(flag LcdStatusFlag) bool {
	return GetBit(lcd.stat, byte(flag))
}

func (lcd *LCD) GetMode() LcdMode {
	return LcdMode(lcd.stat & 0b11)
}

func (lcd *LCD) SetMode(mode LcdMode) {
	lcd.stat = lcd.stat&0b11111100 | byte(mode)
}

func (lcd *LCD) updatePalette(palette *[4]uint32, data byte) {
	palette[0] = DEFAULT_PALLETTE[(data>>0)&0b11]
	palette[1] = DEFAULT_PALLETTE[(data>>2)&0b11]
	palette[2] = DEFAULT_PALLETTE[(data>>4)&0b11]
	palette[3] = DEFAULT_PALLETTE[(data>>6)&0b11]
}

func (lcd *LCD) IncrementLy() {
	lcd.ly++

	if lcd.ly == lcd.lyc {
		// Keep the LYC status flag up to date
		lcd.SetLcdStatusFlag(STAT_LYC_EQUAL, true)

		// If the LCD is requesting LYC interrupts, request the interrupt
		if lcd.CheckLcdStatusFlag(STAT_LYC_INTERRUPT) {
			lcd.gameboy.RequestInterrupt(INT_LCD)
		}
	} else {
		lcd.SetLcdStatusFlag(STAT_LYC_EQUAL, false)
	}
}
