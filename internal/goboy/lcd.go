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
	LCDC_ADDRESS uint16 = 0xFF40
	STAT_ADDRESS uint16 = 0xFF41
	LY_ADDRESS   uint16 = 0xFF44
	LYC_ADDRESS  uint16 = 0xFF45
	SCY_ADDRESS  uint16 = 0xFF42
	SCX_ADDRESS  uint16 = 0xFF43
	WX_ADDRESS   uint16 = 0xFF4A
	WY_ADDRESS   uint16 = 0xFF4B
	BGP_ADDRESS  uint16 = 0xFF47
	OBJ0_ADDRESS uint16 = 0xFF48
	OBJ1_ADDRESS uint16 = 0xFF49
)

var DEFAULT_PALLETTE = [4]uint32{0xFFFFFFFF, 0xFFAAAAAA, 0xFF555555, 0xFF000000}

type LCD struct {
	bus *Bus

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

func NewLCD(bus *Bus) *LCD {
	return &LCD{
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
	case LCDC_ADDRESS:
		return lcd.lcdc
	case STAT_ADDRESS:
		return lcd.stat
	case LY_ADDRESS:
		return lcd.ly
	case LYC_ADDRESS:
		return lcd.lyc
	case SCY_ADDRESS:
		return lcd.scy
	case SCX_ADDRESS:
		return lcd.scx
	case WX_ADDRESS:
		return lcd.wx
	case WY_ADDRESS:
		return lcd.wy
	case BGP_ADDRESS:
		return lcd.bgp
	case OBJ0_ADDRESS:
		return lcd.obj0
	case OBJ1_ADDRESS:
		return lcd.obj1
	}

	panic(fmt.Sprintf("Bad read from LCD: %4.4X", address))
}

func (lcd *LCD) writeByte(address uint16, value byte) {
	switch address {
	case LCDC_ADDRESS:
		lcd.lcdc = value
		return
	case STAT_ADDRESS:
		lcd.stat = value
		return
	case LY_ADDRESS:
		lcd.ly = value
		return
	case LYC_ADDRESS:
		lcd.lyc = value
		return
	case SCY_ADDRESS:
		lcd.scy = value
		return
	case SCX_ADDRESS:
		lcd.scx = value
		return
	case WX_ADDRESS:
		lcd.wx = value
		return
	case WY_ADDRESS:
		lcd.wy = value
		return
	case BGP_ADDRESS:
		lcd.updatePalette(&lcd.bgColors, value)
		lcd.bgp = value
		return
	case OBJ0_ADDRESS:
		lcd.updatePalette(&lcd.sp1Colors, value&0b11111100)
		lcd.obj0 = value
		return
	case OBJ1_ADDRESS:
		lcd.updatePalette(&lcd.sp2Colors, value&0b11111100)
		lcd.obj1 = value
		return
	}

	panic(fmt.Sprintf("Bad write to LCD: %4.4X", address))
}

func (lcd *LCD) SetLcdControlFlag(flag LcdControlFlag, on bool) {
	lcd.lcdc = SetBit(lcd.lcdc, byte(flag), on)
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
	if GetBit(lcd.lcdc, byte(BGW_TILES)) {
		return 0x9C00
	}

	return 0x9800
}

func (lcd *LCD) IsLcdEnabled() bool {
	return GetBit(lcd.lcdc, byte(LCD_ENABLE))
}

func (lcd *LCD) SetLcdStatusFlag(flag LcdStatusFlag, on bool) {
	lcd.lcdc = SetBit(lcd.stat, byte(flag), on)
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
	// Increment LY
	lcd.ly++

	if lcd.ly == lcd.lyc {
		// Keep the LYC status flag up to date
		lcd.SetLcdStatusFlag(STAT_LYC_EQUAL, true)

		// If the LCD is requesting LYC interrupts, request the interrupt
		if lcd.CheckLcdStatusFlag(STAT_LYC_INTERRUPT) {
			lcd.bus.interruptEnableRegister.SetInterrupt(INT_LCD, true)
		}
	} else {
		lcd.SetLcdStatusFlag(STAT_LYC_EQUAL, false)
	}
}
