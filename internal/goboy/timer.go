package goboy

// @see https://gbdev.io/pandocs/Timer_and_Divider_Registers.html
type Timer struct {
	// Reference to the GameBoy this timer is a part of. The timer is capable of
	// requesting interrupts, so it needs access to the main hardware to do so.
	gameboy *GameBoy
	// This register is incremented at a rate of 16384Hz (~16779Hz on SGB).
	// Writing any value to this register resets it to $00. Additionally, this
	// register is reset when executing the stop instruction, and only begins
	// ticking again once stop mode ends.
	div byte
	// This timer is incremented at the clock frequency specified by the TAC
	// register ($FF07). When the value overflows (exceeds $FF) it is reset to the
	// value specified in TMA (FF06) and an interrupt is requested
	tima byte
	// 	When TIMA overflows, it is reset to the value in this register and an
	// 	interrupt is requested. Example of use: if TMA is set to $FF, an interrupt
	// 	is requested at the clock frequency selected in TAC (because every increment
	// 	is an overflow). However, if TMA is set to $FE, an interrupt is only
	// 	requested every two increments, which effectively divides the selected clock
	// 	by two. Setting TMA to $FD would divide the clock by three, and so on.
	// If a TMA write is executed on the same M-cycle as the content of TMA is
	// transferred to TIMA due to a timer overflow, the old value is transferred to
	// TIMA.
	tma byte
	tac byte
}

func NewTimer(gameboy *GameBoy) *Timer {
	return &Timer{
		gameboy: gameboy,
		div:     0x1E,
		tima:    0,
		tma:     0,
		tac:     0xF8,
	}
}

func (timer *Timer) Tick(cpu *CPU) {
	lastDiv := timer.div
	timer.div++

	timerUpdate := false

	switch timer.tac & (0b11) {
	case 0b00:
		timerUpdate = lastDiv > timer.div
	case 0b01:
		timerUpdate = GetBit(lastDiv, 3) && !GetBit(timer.div, 3)
	case 0b10:
		timerUpdate = GetBit(lastDiv, 5) && !GetBit(timer.div, 5)
	case 0b11:
		timerUpdate = GetBit(lastDiv, 7) && !GetBit(timer.div, 7)
	}

	if timerUpdate && GetBit(timer.tac, 2) {
		timer.tima++

		if timer.tima == 0xFF {
			timer.tima = timer.tma
			timer.gameboy.RequestInterrupt(INT_TIMER)
		}
	}
}

func (timer *Timer) readByte(address uint16) byte {
	switch address {
	case 0xFF04:
		return timer.div
	case 0xFF05:
		return timer.tima
	case 0xFF06:
		return timer.tma
	case 0xFF07:
		return timer.tac
	}

	panic("Attempted to read invalid timer address")
}

func (timer *Timer) writeByte(address uint16, value byte) {
	switch address {
	case 0xFF04:
		timer.div = 0
	case 0xFF05:
		timer.tima = value
	case 0xFF06:
		timer.tma = value
	case 0xFF07:
		timer.tac = value
	}
}
