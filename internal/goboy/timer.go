package goboy

const (
	TIMER_DIV  = 0xFF04
	TIMER_TIMA = 0xFF05
	TIMER_TMA  = 0xFF06
	TIMER_TAC  = 0xFF07
)

// @see https://gbdev.io/pandocs/Timer_and_Divider_Registers.html
type Timer struct {
	// Reference to the GameBoy this timer is a part of. The timer is capable of
	// requesting interrupts, so it needs access to the main hardware to do so.
	gameboy *GameBoy
	// This register is incremented at a rate of 16384Hz (~16779Hz on SGB).
	// Writing any value to this register resets it to $00. Additionally, this
	// register is reset when executing the stop instruction, and only begins
	// ticking again once stop mode ends. DIV is the upper 8 bits of an internal
	// 16 bit register which
	sysclk uint16
	div    byte
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
	// Controls whether or not the timer is enabled, and how fast it ticks at if
	// it is enabled.
	tac byte

	timerBit byte
}

// Map of TAC lo bits to the bit that needs to roll over for TIMA to be
// incremented
// @see https://gbdev.io/pandocs/Timer_and_Divider_Registers.html#ff07--tac-timer-control
var timerBitMap = [4]byte{9, 3, 5, 7}

func NewTimer(gameboy *GameBoy) *Timer {
	return &Timer{
		gameboy: gameboy,
		sysclk:  0x1E00,
		div:     0x1E,
		tima:    0,
		tma:     0,
		tac:     0xF8,

		timerBit: 9,
	}
}

func (timer *Timer) Tick() {
	lastSysclk := timer.sysclk
	timer.sysclk += 1
	timer.div = byte(timer.sysclk >> 8)

	if !timer.enabled() {
		return
	}

	incrementTima := (lastSysclk&(1<<timer.timerBit)) != 0 && (timer.sysclk&(1<<timer.timerBit)) == 0

	if !incrementTima {
		return
	}

	timer.tima++

	// When TIMA overflows it should be reset back to the value of TMA and an
	// interrupt should be requested
	if timer.tima == 0x00 {
		timer.tima = timer.tma
		timer.gameboy.RequestInterrupt(INT_TIMER)
	}
}

func (timer *Timer) enabled() bool {
	return GetBit(timer.tac, 2)
}

func (timer *Timer) readByte(address uint16) byte {
	switch address {
	case TIMER_DIV:
		return timer.div
	case TIMER_TIMA:
		return timer.tima
	case TIMER_TMA:
		return timer.tma
	case TIMER_TAC:
		return timer.tac
	}

	panic("Attempted to read invalid timer address")
}

func (timer *Timer) writeByte(address uint16, value byte) {
	switch address {
	case TIMER_DIV:
		timer.sysclk = 0
		timer.div = 0
	case TIMER_TIMA:
		timer.tima = value
	case TIMER_TMA:
		timer.tma = value
	case TIMER_TAC:
		timer.tac = value
		timer.timerBit = timerBitMap[value&0b11]
	}
}
