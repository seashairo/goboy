package goboy

type Button byte

const (
	JOYPAD_A = Button(iota)
	JOYPAD_B
	JOYPAD_SELECT
	JOYPAD_START
	JOYPAD_RIGHT
	JOYPAD_LEFT
	JOYPAD_UP
	JOYPAD_DOWN
)

const JOYPAD_SELECT_D_PAD_BIT = 4
const JOYPAD_SELECT_BUTTONS_BIT = 5

// @see https://gbdev.io/pandocs/Joypad_Input.html#ff00--p1joyp-joypad
type Joypad struct {
	bus     *Bus
	data    byte
	buttons byte
}

func NewJoypad(bus *Bus) *Joypad {
	return &Joypad{
		bus:     bus,
		data:    0xFF,
		buttons: 00,
	}
}

func (joypad *Joypad) writeByte(_ uint16, value byte) {
	// the lower nibble of the joypad is read-only
	joypad.data = (joypad.data & 0x0F) | (value & 0xF0)
}

func (joypad *Joypad) readByte(_ uint16) byte {
	dpad := !GetBit(joypad.data, JOYPAD_SELECT_D_PAD_BIT)
	buttons := !GetBit(joypad.data, JOYPAD_SELECT_BUTTONS_BIT)

	hi := joypad.data & 0xF0
	lo := byte(0x0F)

	if buttons && !dpad {
		lo = byte(^joypad.buttons) & 0x0F
	} else if dpad && !buttons {
		lo = (byte(^joypad.buttons) & 0xF0) >> 4
	}
	return hi | lo
}

func (joypad *Joypad) Press(button Button) {
	if !joypad.Check(button) {
		joypad.buttons = SetBit(joypad.buttons, byte(button), true)
	}
}

func (joypad *Joypad) Release(button Button) {
	if joypad.Check(button) {
		joypad.buttons = SetBit(joypad.buttons, byte(button), false)
	}
}

func (joypad *Joypad) Check(button Button) bool {
	return GetBit(joypad.buttons, byte(button))
}
