package goboy

type LengthCounter struct {
	enabled        bool
	length         byte
	fullLength     byte
	frameSequencer byte
}

func NewLengthCounter() *LengthCounter {
	return &LengthCounter{
		enabled:        false,
		length:         0,
		fullLength:     0,
		frameSequencer: 0,
	}
}

func (lc *LengthCounter) Tick() {
	if lc.enabled && lc.length > 0 {
		lc.length--
	}
}

func (lc *LengthCounter) SetNR4(value byte) {
	enable := GetBit(value, 6)
	trigger := GetBit(value, 7)

	if lc.enabled {
		if trigger && lc.length == 0 {
			if enable && lc.frameSequencer&1 != 0 {
				lc.length = lc.fullLength - 1
			} else {
				lc.length = lc.fullLength
			}
		}
	} else if enable {
		if lc.frameSequencer&1 != 0 {
			if lc.length != 0 {
				lc.length--
			}
			if trigger && lc.length == 0 {
				lc.length = lc.fullLength - 1
			}
		}
	} else {
		if trigger && lc.length == 0 {
			lc.length = lc.fullLength
		}
	}

	lc.enabled = enable
}

func (lc *LengthCounter) IsEnabled() bool {
	return lc.enabled
}

func (lc *LengthCounter) IsZero() bool {
	return lc.length == 0
}

func (lc *LengthCounter) SetLength(length byte) {
	if length == 0 {
		lc.length = lc.fullLength
	} else {
		lc.length = lc.fullLength - length
	}
}

func (lc *LengthCounter) SetFullLength(fullLength byte) {
	lc.fullLength = fullLength
}

func (lc *LengthCounter) PowerOff() {
	lc.enabled = false
	lc.frameSequencer = 0
}

func (lc *LengthCounter) SetFrameSequencer(frameSequencer byte) {
	lc.frameSequencer = frameSequencer
}
