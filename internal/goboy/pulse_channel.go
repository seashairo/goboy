package goboy

var dutyTable = [4][8]bool{
	{false, false, false, false, false, false, false, true},
	{true, false, false, false, false, false, false, true},
	{true, false, false, false, false, true, true, true},
	{false, true, true, true, true, true, true, false},
}

type PulseChannel struct {
	offset uint16
	onR    bool
	onL    bool

	enabled    bool
	dacEnabled bool

	timer     uint16
	sequence  byte
	duty      byte
	frequency uint16
	output    int16

	volumeEnvelope *VolumeEnvelope
	lengthCounter  *LengthCounter
}

func NewPulseChannel(offset uint16) *PulseChannel {
	return &PulseChannel{
		offset: offset,

		enabled:    false,
		dacEnabled: false,

		duty:      0,
		frequency: 0,
		output:    0,

		volumeEnvelope: NewVolumeEnvelope(),
		lengthCounter:  NewLengthCounter(),
	}
}

func (pc *PulseChannel) GetSample() (int16, int16) {
	left := int16(0)
	if pc.onL {
		left = pc.output
	}

	right := int16(0)
	if pc.onR {
		right = pc.output
	}

	return left, right
}

func (pc *PulseChannel) Tick() {
	pc.timer--
	if pc.timer <= 0 {
		pc.timer = (2048 - pc.frequency) * 4
		pc.sequence = (pc.sequence + 1) & 7

		if pc.enabled && pc.dacEnabled && dutyTable[pc.duty][pc.sequence] {
			pc.output = int16(pc.volumeEnvelope.GetVolume()) * 100
		} else {
			pc.output = 0
		}
	}
}

func (pc *PulseChannel) trigger() {
	pc.timer = (2048 - pc.frequency) * 4
	pc.volumeEnvelope.Trigger()
	pc.enabled = pc.dacEnabled
}

func (pc *PulseChannel) readByte(address uint16) byte {
	switch address - pc.offset {
	case 0:
		return 0xFF
	case 1:
		return pc.duty << 6
	case 2:
		return pc.volumeEnvelope.GetNR2()
	case 3:
		return 0xFF
	case 4:
		return SetBit(0b10111111, 6, pc.lengthCounter.IsEnabled())
	}

	return 0
}

func (pc *PulseChannel) writeByte(address uint16, value byte) {
	switch address - pc.offset {
	case 0:
		return
	case 1:
		pc.duty = value >> 6
		pc.lengthCounter.SetLength(value & 0b00111111)
		return
	case 2:
		pc.dacEnabled = value&0b11111000 != 0
		pc.enabled = pc.enabled && pc.dacEnabled

		pc.volumeEnvelope.SetNR2(value)
		return
	case 3:
		pc.frequency = (pc.frequency & 0b00000111_00000000) | uint16(value)
		return
	case 4:
		pc.frequency = (pc.frequency & 0xFF) | (uint16(value&0b0000111) << 8)
		pc.lengthCounter.SetNR4(value)

		if pc.lengthCounter.IsEnabled() && pc.lengthCounter.IsZero() {
			pc.enabled = false
		} else if GetBit(value, 7) {
			pc.trigger()
		}

		return
	}
}

func (pc *PulseChannel) LengthClock() {
	pc.lengthCounter.Tick()

	if pc.lengthCounter.IsEnabled() && pc.lengthCounter.IsZero() {
		pc.enabled = false
	}
}

func (pc *PulseChannel) EnvelopeClock() {
	pc.volumeEnvelope.Tick()
}

func (pc *PulseChannel) PowerOff() {
	pc.volumeEnvelope.PowerOff()
	pc.lengthCounter.PowerOff()

	pc.enabled = false
	pc.dacEnabled = false
	pc.onL = false
	pc.onR = false

	pc.sequence = 0
	pc.frequency = 0
	pc.duty = 0
	pc.output = 0
}
