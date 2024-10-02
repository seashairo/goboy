package goboy

// @see https://gbdev.io/pandocs/Audio_Registers.html
const (
	// Sound Channel 1 — Pulse with period sweep
	APU_NR10 = 0xFF10 // sweep
	APU_NR11 = 0xFF11 // length timer & duty cycle
	APU_NR12 = 0xFF12 // volume & envelope
	APU_NR13 = 0xFF13 // period low [write-only]
	APU_NR14 = 0xFF14 // period high & control

	// Sound Channel 2 — Pulse
	APU_NR20 = 0xFF15 // unused
	APU_NR21 = 0xFF16 // length timer & duty cycle
	APU_NR22 = 0xFF17 // volume & envelope
	APU_NR23 = 0xFF18 // period low [write-only]
	APU_NR24 = 0xFF19 // period high & control

	// Sound Channel 3 — Wave output
	APU_NR30 = 0xFF1A // DAC enable
	APU_NR31 = 0xFF1B // length timer [write-only]
	APU_NR32 = 0xFF1C // output level
	APU_NR33 = 0xFF1D // period low [write-only]
	APU_NR34 = 0xFF1E // period high & control

	// Sound Channel 4 — Noise
	APU_NR41 = 0xFF20 // length timer [write-only]
	APU_NR42 = 0xFF21 // volume & envelope
	APU_NR43 = 0xFF22 // frequency & randomness
	APU_NR44 = 0xFF23 // control

	// Global control registers
	APU_NR50 = 0xFF24 // Master volume & VIN panning
	APU_NR51 = 0xFF25 // Sound panning
	APU_NR52 = 0xFF26 // Audio master control

	// Wave pattern RAM
	APU_WAVE_RAM_START = 0xFF30
	APU_WAVE_RAM_END   = 0xFF3F
)

const (
	sampleRate        = 44100
	cpuTicksPerSample = 4194304 / sampleRate
)

type AudioCallback func(left int16, right int16)

// @see https://gbdev.io/pandocs/Audio.html
type APU struct {
	gameboy *GameBoy

	enabled bool

	nr11, nr12, nr13, nr14       byte
	nr30, nr31, nr32, nr33, nr34 byte
	nr41, nr42, nr43, nr44       byte

	waveRam *RAM
	// pulseChannel1 *AudioChannel
	pulseChannel2 *PulseChannel
	// waveChannel   *AudioChannel
	// noiseChannel  *AudioChannel

	onL, onR         bool
	volumeL, volumeR byte

	// Internal clock counters
	frameSequencerCounter uint16
	frameSequencer        byte
	sampleCounter         byte

	callbacks []AudioCallback
}

func NewAPU(gameboy *GameBoy) *APU {
	apu := &APU{
		gameboy:   gameboy,
		waveRam:   NewRAM(16, APU_WAVE_RAM_START),
		callbacks: make([]AudioCallback, 0),

		frameSequencerCounter: 8192,
		frameSequencer:        0,
		sampleCounter:         cpuTicksPerSample,
	}

	// apu.pulseChannel1 = NewAudioChannel()
	apu.pulseChannel2 = NewPulseChannel(APU_NR20 - 1)
	// apu.waveChannel = NewAudioChannel()
	// apu.noiseChannel = NewAudioChannel()

	return apu
}

func (apu *APU) Tick() {
	if !apu.enabled {
		return
	}

	apu.frameSequencerCounter -= 1
	if apu.frameSequencerCounter == 0 {
		apu.frameSequencerCounter = 8192

		switch apu.frameSequencer {
		case 0:
			apu.pulseChannel2.LengthClock()
		case 2:
			apu.pulseChannel2.LengthClock()
		case 4:
			apu.pulseChannel2.LengthClock()
		case 6:
			apu.pulseChannel2.LengthClock()
		case 7:
			apu.pulseChannel2.EnvelopeClock()
		}

		apu.frameSequencer = (apu.frameSequencer + 1) & 7

		apu.pulseChannel2.lengthCounter.SetFrameSequencer(apu.frameSequencer)
	}

	apu.pulseChannel2.Tick()

	apu.sampleCounter -= 1
	if apu.sampleCounter != 0 {
		return
	}

	apu.sampleCounter = cpuTicksPerSample

	pulseSampleL, pulseSampleR := apu.pulseChannel2.GetSample()

	finalSampleL := pulseSampleL
	finalSampleR := pulseSampleR

	for _, cb := range apu.callbacks {
		if cb != nil {
			cb(finalSampleL, finalSampleR)
		}
	}
}

func (apu *APU) readByte(address uint16) byte {
	if Between(address, APU_NR20, APU_NR24) {
		return apu.pulseChannel2.readByte(address)
	}

	switch address {
	case APU_NR10:
		return 0
	case APU_NR11:
		// the lo 6 bits of nr11 are write-only
		return apu.nr11 & 0b11000000
	case APU_NR12:
		return apu.nr12
	case APU_NR13:
		return 0x00
	case APU_NR14:
		return apu.nr14
	case APU_NR30:
		return apu.nr30
	case APU_NR31:
		return apu.nr31
	case APU_NR32:
		return apu.nr32
	case APU_NR33:
		return apu.nr33
	case APU_NR34:
		return apu.nr34
	case APU_NR41:
		return apu.nr41
	case APU_NR42:
		return apu.nr42
	case APU_NR43:
		return apu.nr43
	case APU_NR44:
		return apu.nr44
	case APU_NR50:
		out := byte(0)

		out = SetBit(out, 7, apu.onL)
		out = SetBit(out, 3, apu.onL)

		out |= (apu.volumeL - 1) << 4
		out |= (apu.volumeR - 1)

		return out
	case APU_NR51:
		out := byte(0)

		// out = SetBit(out, 0, apu.pulseChannel1.onR)
		out = SetBit(out, 1, apu.pulseChannel2.onR)
		// out = SetBit(out, 2, apu.waveChannel.onR)
		// out = SetBit(out, 3, apu.noiseChannel.onR)
		// out = SetBit(out, 4, apu.pulseChannel1.onL)
		out = SetBit(out, 5, apu.pulseChannel2.onL)
		// out = SetBit(out, 6, apu.waveChannel.onL)
		// out = SetBit(out, 7, apu.noiseChannel.onL)

		return out
	case APU_NR52:
		out := SetBit(0, 7, apu.enabled)
		out = SetBit(out, 1, apu.pulseChannel2.enabled)
		return out
	}

	if Between(address, APU_WAVE_RAM_START, APU_WAVE_RAM_END) {
		return apu.waveRam.readByte(address)
	}

	return 0
}

func (apu *APU) writeByte(address uint16, value byte) {
	// If audio master is not enabled, then the APU is considered read-only with
	// the exception of NR52 which controls whether the APU is enabled.
	if !apu.enabled && address != APU_NR52 {
		return
	}

	if Between(address, APU_NR20, APU_NR24) {
		apu.pulseChannel2.writeByte(address, value)
		return
	}

	switch address {
	case APU_NR10:
		return
	case APU_NR11:
		apu.nr11 = value
		return
	case APU_NR12:
		apu.nr12 = value
		return
	case APU_NR13:
		return
	case APU_NR14:
		return
	case APU_NR30:
		apu.nr30 = value
		return
	case APU_NR31:
		apu.nr31 = value
		return
	case APU_NR32:
		apu.nr32 = value
		return
	case APU_NR33:
		apu.nr33 = value
		return
	case APU_NR34:
		apu.nr34 = value
		return
	case APU_NR41:
		apu.nr41 = value
		return
	case APU_NR42:
		apu.nr42 = value
		return
	case APU_NR43:
		apu.nr43 = value
		return
	case APU_NR44:
		apu.nr44 = value
		return
	case APU_NR50:
		apu.onL = GetBit(value, 7)
		apu.volumeL = (value & 0b01110000) + 1
		apu.onR = GetBit(value, 3)
		apu.volumeR = (value & 0b00000111) + 1
		return
	case APU_NR51:
		// apu.pulseChannel1.onR = GetBit(value, 0)
		apu.pulseChannel2.onR = GetBit(value, 1)
		// apu.waveChannel.onR = GetBit(value, 2)
		// apu.noiseChannel.onR = GetBit(value, 3)
		// apu.pulseChannel1.onL = GetBit(value, 4)
		apu.pulseChannel2.onL = GetBit(value, 5)
		// apu.waveChannel.onL = GetBit(value, 6)
		// apu.noiseChannel.onL = GetBit(value, 7)
		return
	case APU_NR52:
		enable := GetBit(value, 7)

		if apu.enabled && !enable {
			apu.enabled = false
			apu.volumeL = 0
			apu.volumeR = 0
			apu.pulseChannel2.PowerOff()
		} else if !apu.enabled && enable {
			apu.frameSequencer = 0
		}

		apu.enabled = enable
		return
	}

	if Between(address, APU_WAVE_RAM_START, APU_WAVE_RAM_END) {
		apu.waveRam.writeByte(address, value)
	}
}
