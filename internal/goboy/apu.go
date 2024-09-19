package goboy

import "math"

// @see https://gbdev.io/pandocs/Audio_Registers.html
const (
	// Sound Channel 1 — Pulse with period sweep
	APU_NR10 = 0xFF10 // sweep
	APU_NR11 = 0xFF11 // length timer & duty cycle
	APU_NR12 = 0xFF12 // volume & envelope
	APU_NR13 = 0xFF13 // period low [write-only]
	APU_NR14 = 0xFF14 // period high & control

	// Sound Channel 2 — Pulse
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

type AudioCallback func(sample int16)

// @see https://gbdev.io/pandocs/Audio.html
type APU struct {
	gameboy *GameBoy

	nr10 byte
	nr11 byte
	nr12 byte
	nr13 byte
	nr14 byte

	nr21 byte
	nr22 byte
	nr23 byte
	nr24 byte

	nr30 byte
	nr31 byte
	nr32 byte
	nr33 byte
	nr34 byte

	nr41 byte
	nr42 byte
	nr43 byte
	nr44 byte

	nr50 byte
	nr51 byte
	nr52 byte

	waveRam *RAM

	callbacks []AudioCallback
}

func NewAPU(gameboy *GameBoy) *APU {
	return &APU{
		gameboy: gameboy,

		nr10: 0,
		nr11: 0,
		nr12: 0,
		nr13: 0,
		nr14: 0,

		nr21: 0,
		nr22: 0,
		nr23: 0,
		nr24: 0,

		nr30: 0,
		nr31: 0,
		nr32: 0,
		nr33: 0,
		nr34: 0,

		nr41: 0,
		nr42: 0,
		nr43: 0,
		nr44: 0,

		nr50: 0,
		nr51: 0,
		nr52: 0,

		waveRam: NewRAM(16, APU_WAVE_RAM_START),

		callbacks: make([]AudioCallback, 0),
	}
}

var phase float64

const (
	amplitude = 1000 // Amplitude of the waveform
	frequency = 440  // Frequency of the sine wave (A4)
)

func (apu *APU) generateSample() int16 {
	sample := int16(amplitude * math.Sin(phase))

	phase += 2 * math.Pi * frequency / float64(sampleRate)
	if phase > 2*math.Pi {
		phase -= 2 * math.Pi
	}

	return sample
}

func (apu *APU) Tick() {
	sample := apu.generateSample()
	for _, cb := range apu.callbacks {
		if cb != nil {
			cb(sample)
		}
	}
}

func (apu *APU) readByte(address uint16) byte {
	switch address {
	case APU_NR10:
		return apu.nr10
	case APU_NR11:
		return apu.nr11
	case APU_NR12:
		return apu.nr12
	case APU_NR13:
		return apu.nr13
	case APU_NR14:
		return apu.nr14
	case APU_NR21:
		return apu.nr21
	case APU_NR22:
		return apu.nr22
	case APU_NR23:
		return apu.nr23
	case APU_NR24:
		return apu.nr24
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
		return apu.nr50
	case APU_NR51:
		return apu.nr51
	case APU_NR52:
		return apu.nr52
	}

	if Between(address, APU_WAVE_RAM_START, APU_WAVE_RAM_END) {
		return apu.waveRam.readByte(address)
	}

	return 0
}

func (apu *APU) writeByte(address uint16, value byte) {
	switch address {
	case APU_NR10:
		apu.nr10 = value
		return
	case APU_NR11:
		apu.nr11 = value
		return
	case APU_NR12:
		apu.nr12 = value
		return
	case APU_NR13:
		apu.nr13 = value
		return
	case APU_NR14:
		apu.nr14 = value
		return
	case APU_NR21:
		apu.nr21 = value
		return
	case APU_NR22:
		apu.nr22 = value
		return
	case APU_NR23:
		apu.nr23 = value
		return
	case APU_NR24:
		apu.nr24 = value
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
		apu.nr50 = value
		return
	case APU_NR51:
		apu.nr51 = value
		return
	case APU_NR52:
		apu.nr52 = value
		return
	}

	if Between(address, APU_WAVE_RAM_START, APU_WAVE_RAM_END) {
		apu.waveRam.writeByte(address, value)
	}
}
