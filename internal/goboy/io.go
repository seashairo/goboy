package goboy

const (
	IO_IF = 0xFF0F
)

type IO struct {
	interrupts *InterruptRegister
	timer      *Timer
	dma        *DMA
	lcd        *LCD
	joypad     *Joypad
	serial     *Serial
	apu        *APU
}

func NewIO(gameboy *GameBoy, bus *Bus, timer *Timer, interruptEnableRegister *InterruptRegister, lcd *LCD, joypad *Joypad, apu *APU) *IO {
	dma := NewDMA(bus)
	serial := NewSerial(gameboy)

	return &IO{
		interrupts: interruptEnableRegister,
		timer:      timer,
		dma:        dma,
		lcd:        lcd,
		joypad:     joypad,
		serial:     serial,
		apu:        apu,
	}
}

func (io *IO) writeByte(address uint16, value byte) {
	if address == IO_JOYP {
		io.joypad.writeByte(address, value)
		return
	}

	if Between(address, SERIAL_SB, SERIAL_SC) {
		io.serial.writeByte(address, value)
		return
	}

	if Between(address, TIMER_DIV, TIMER_TAC) {
		io.timer.writeByte(address, value)
		return
	}

	if Between(address, APU_NR10, APU_WAVE_RAM_END) {
		io.apu.writeByte(address, value)
		return
	}

	if address == IO_IF {
		io.interrupts.writeByte(value)
		return
	}

	if address == IO_DMA {
		io.dma.writeByte(address, value)
		return
	}

	if Between(address, LCD_LCDC, LCD_WX) {
		io.lcd.writeByte(address, value)
		return
	}

	// fmt.Printf("Writing to %2.2X not supported (IO_REGISTERS)\n", address)
}

func (io *IO) readByte(address uint16) byte {
	if address == IO_JOYP {
		return io.joypad.readByte(address)
	}

	if Between(address, SERIAL_SB, SERIAL_SC) {
		return io.serial.readByte(address)
	}

	if Between(address, TIMER_DIV, TIMER_TAC) {
		return io.timer.readByte(address)
	}

	if Between(address, APU_NR10, APU_WAVE_RAM_END) {
		return io.apu.readByte(address)
	}

	if address == IO_IF {
		return io.interrupts.readByte()
	}

	if address == IO_DMA {
		return io.dma.readByte(address)
	}

	if Between(address, LCD_LCDC, LCD_WX) {
		return io.lcd.readByte(address)
	}

	// fmt.Printf("Reading from %2.2X not supported (IO_REGISTERS)\n", address)
	return 0
}
