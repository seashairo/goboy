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
	sampleRate      = 44100
	clocksPerSecond = 4194304
	clocksPerSample = clocksPerSecond / sampleRate / 2
	clocksPerFrame  = 8192
)

type AudioCallback func(left int16, right int16)

type APU struct {
	gameboy *GameBoy

	LeftSample  uint32
	RightSample uint32
	NumSamples  uint32

	LastLeft           float64
	LastRight          float64
	LastCorrectedLeft  float64
	LastCorrectedRight float64

	masterEnable bool

	frameSequencerCounter uint16
	frameSequencer        byte

	soundChannels [4]SoundChannel

	VInToLeftSpeaker  bool
	VInToRightSpeaker bool

	RightSpeakerVolume byte
	LeftSpeakerVolume  byte

	callbacks []AudioCallback
}

func NewAPU(gameboy *GameBoy) *APU {
	apu := &APU{
		gameboy:               gameboy,
		frameSequencerCounter: clocksPerFrame,
	}

	apu.soundChannels[0].soundType = squareSoundType
	apu.soundChannels[1].soundType = squareSoundType
	apu.soundChannels[2].soundType = waveSoundType
	apu.soundChannels[3].soundType = noiseSoundType

	apu.soundChannels[3].polyFeedbackReg = 0x01

	return apu
}

func (apu *APU) generateSample() {
	apu.TickFrequency()

	leftSam, rightSam := uint32(0), uint32(0)
	if apu.masterEnable {
		left0, right0 := apu.soundChannels[0].getSample()
		leftSam += uint32(left0)
		rightSam += uint32(right0)

		left1, right1 := apu.soundChannels[1].getSample()
		leftSam += uint32(left1)
		rightSam += uint32(right1)

		left2, right2 := apu.soundChannels[2].getSample()
		leftSam += uint32(left2)
		rightSam += uint32(right2)

		left3, right3 := apu.soundChannels[3].getSample()
		leftSam += uint32(left3)
		rightSam += uint32(right3)

		leftSam *= uint32(apu.LeftSpeakerVolume + 1)
		rightSam *= uint32(apu.RightSpeakerVolume + 1)
	}

	apu.LeftSample += leftSam
	apu.RightSample += rightSam
	apu.NumSamples++

	if apu.NumSamples >= clocksPerSample {
		left := float64(apu.LeftSample) / float64(apu.NumSamples)
		right := float64(apu.RightSample) / float64(apu.NumSamples)
		left /= 4 * 8 * 15
		right /= 4 * 8 * 15

		correctedLeft := left - apu.LastLeft + 0.995*apu.LastCorrectedLeft
		apu.LastCorrectedLeft = correctedLeft
		apu.LastLeft = left
		left = correctedLeft

		correctedRight := right - apu.LastRight + 0.995*apu.LastCorrectedRight
		apu.LastCorrectedRight = correctedRight
		apu.LastRight = right
		right = correctedRight

		for _, cb := range apu.callbacks {
			if cb != nil {
				cb(int16(left*32767.0), int16(right*32767.0))
			}
		}

		apu.LeftSample = 0
		apu.RightSample = 0
		apu.NumSamples = 0
	}
}

func (apu *APU) Tick() {
	apu.frameSequencerCounter--
	if apu.frameSequencerCounter == 0 {
		apu.frameSequencerCounter = clocksPerFrame

		switch apu.frameSequencer {
		case 0:
			apu.TickLength()
		case 2:
			apu.TickLength()
			apu.soundChannels[0].TickSweep()
		case 4:
			apu.TickLength()
		case 6:
			apu.TickLength()
			apu.soundChannels[0].TickSweep()
		case 7:
			apu.TickVolumeEnvelope()
		}

		apu.frameSequencer = (apu.frameSequencer + 1) & 7
	}

	if apu.frameSequencerCounter&1 == 0 {
		apu.generateSample()
	}
}

func (apu *APU) TickFrequency() {
	apu.soundChannels[0].TickFrequency()
	apu.soundChannels[1].TickFrequency()
	apu.soundChannels[2].TickFrequency()
	apu.soundChannels[3].TickFrequency()
}

func (apu *APU) TickLength() {
	apu.soundChannels[0].TickLength()
	apu.soundChannels[1].TickLength()
	apu.soundChannels[2].TickLength()
	apu.soundChannels[3].TickLength()
}

func (apu *APU) TickVolumeEnvelope() {
	apu.soundChannels[0].TickVolumeEnvelope()
	apu.soundChannels[1].TickVolumeEnvelope()
	apu.soundChannels[2].TickVolumeEnvelope()
	apu.soundChannels[3].TickVolumeEnvelope()
}

func (apu *APU) writeByte(address uint16, value byte) {
	switch address {
	case APU_NR10:
		apu.soundChannels[0].writeSweepReg(value)
	case APU_NR11:
		apu.soundChannels[0].writeLenDutyReg(value)
	case APU_NR12:
		apu.soundChannels[0].writeVolumeEnvelope(value)
	case APU_NR13:
		apu.soundChannels[0].writePeriodLow(value)
	case APU_NR14:
		apu.soundChannels[0].writePeriodHigh(value)

	case APU_NR21:
		apu.soundChannels[1].writeLenDutyReg(value)
	case APU_NR22:
		apu.soundChannels[1].writeVolumeEnvelope(value)
	case APU_NR23:
		apu.soundChannels[1].writePeriodLow(value)
	case APU_NR24:
		apu.soundChannels[1].writePeriodHigh(value)

	case APU_NR30:
		apu.soundChannels[2].writeWaveOnOffReg(value)
	case APU_NR31:
		apu.soundChannels[2].writeLengthDataReg(value)
	case APU_NR32:
		apu.soundChannels[2].writeWaveOutLvlReg(value)
	case APU_NR33:
		apu.soundChannels[2].writePeriodLow(value)
	case APU_NR34:
		apu.soundChannels[2].writePeriodHigh(value)

	case APU_NR41:
		apu.soundChannels[3].writeLengthDataReg(value)
	case APU_NR42:
		apu.soundChannels[3].writeVolumeEnvelope(value)
	case APU_NR43:
		apu.soundChannels[3].writePolyCounterReg(value)
	case APU_NR44:
		apu.soundChannels[3].writePeriodHigh(value)

	case APU_NR50:
		apu.writeVolumeReg(value)
	case APU_NR51:
		apu.writeSpeakerSelectReg(value)
	case APU_NR52:
		apu.writeSoundOnOffReg(value)
	}

	if address >= APU_WAVE_RAM_START && address <= APU_WAVE_RAM_END {
		apu.soundChannels[2].writeWavePatternValue(address-APU_WAVE_RAM_START, value)
	}
}

func (apu *APU) readByte(address uint16) byte {
	switch address {
	case APU_NR10:
		return apu.soundChannels[0].readSweepReg()
	case APU_NR11:
		return apu.soundChannels[0].readLenDutyReg()
	case APU_NR12:
		return apu.soundChannels[0].readVolumeEnvelope()
	case APU_NR13:
		return apu.soundChannels[0].readPeriodLow()
	case APU_NR14:
		return apu.soundChannels[0].readPeriodHigh()

	case APU_NR21:
		return apu.soundChannels[1].readLenDutyReg()
	case APU_NR22:
		return apu.soundChannels[1].readVolumeEnvelope()
	case APU_NR23:
		return apu.soundChannels[1].readPeriodLow()
	case APU_NR24:
		return apu.soundChannels[1].readPeriodHigh()

	case APU_NR30:
		return apu.soundChannels[2].readWaveOnOffReg()
	case APU_NR31:
		return apu.soundChannels[2].readLengthDataReg()
	case APU_NR32:
		return apu.soundChannels[2].readWaveOutLvlReg()
	case APU_NR33:
		return apu.soundChannels[2].readPeriodLow()
	case APU_NR34:
		return apu.soundChannels[2].readPeriodHigh()

	case APU_NR41:
		return apu.soundChannels[3].readLengthDataReg()
	case APU_NR42:
		return apu.soundChannels[3].readVolumeEnvelope()
	case APU_NR43:
		return apu.soundChannels[3].readPolyCounterReg()
	case APU_NR44:
		return apu.soundChannels[3].readPeriodHigh()

	case APU_NR50:
		return apu.readVolumeReg()
	case APU_NR51:
		return apu.readSpeakerSelectReg()
	case APU_NR52:
		return apu.readSoundOnOffReg()
	}

	if address >= APU_WAVE_RAM_START && address < APU_WAVE_RAM_END {
		return apu.soundChannels[2].wavePatternRAM[address-APU_WAVE_RAM_START]
	}

	return 0xFF
}

func (apu *APU) writeVolumeReg(value byte) {
	apu.VInToLeftSpeaker = GetBit(value, 7)
	apu.VInToRightSpeaker = GetBit(value, 3)
	apu.RightSpeakerVolume = (value >> 4) & 0x07
	apu.LeftSpeakerVolume = value & 0x07
}

func (apu *APU) readVolumeReg() byte {
	out := apu.RightSpeakerVolume<<4 | apu.LeftSpeakerVolume
	out = SetBit(out, 7, apu.VInToLeftSpeaker)
	out = SetBit(out, 3, apu.VInToRightSpeaker)
	return out
}

func (apu *APU) writeSpeakerSelectReg(value byte) {
	apu.soundChannels[0].rightSpeakerOn = GetBit(value, 0)
	apu.soundChannels[1].rightSpeakerOn = GetBit(value, 1)
	apu.soundChannels[2].rightSpeakerOn = GetBit(value, 2)
	apu.soundChannels[3].rightSpeakerOn = GetBit(value, 3)
	apu.soundChannels[0].leftSpeakerOn = GetBit(value, 4)
	apu.soundChannels[1].leftSpeakerOn = GetBit(value, 5)
	apu.soundChannels[2].leftSpeakerOn = GetBit(value, 6)
	apu.soundChannels[3].leftSpeakerOn = GetBit(value, 7)
}

func (apu *APU) readSpeakerSelectReg() byte {
	out := byte(0)

	out = SetBit(out, 0, apu.soundChannels[0].rightSpeakerOn)
	out = SetBit(out, 1, apu.soundChannels[1].rightSpeakerOn)
	out = SetBit(out, 2, apu.soundChannels[2].rightSpeakerOn)
	out = SetBit(out, 3, apu.soundChannels[3].rightSpeakerOn)
	out = SetBit(out, 4, apu.soundChannels[0].leftSpeakerOn)
	out = SetBit(out, 5, apu.soundChannels[1].leftSpeakerOn)
	out = SetBit(out, 6, apu.soundChannels[2].leftSpeakerOn)
	out = SetBit(out, 7, apu.soundChannels[3].leftSpeakerOn)

	return out
}

func (apu *APU) writeSoundOnOffReg(value byte) {
	apu.masterEnable = GetBit(value, 7)
}

func (apu *APU) readSoundOnOffReg() byte {
	out := byte(0b01110000)

	out = SetBit(out, 0, apu.soundChannels[0].enabled)
	out = SetBit(out, 1, apu.soundChannels[1].enabled)
	out = SetBit(out, 2, apu.soundChannels[2].enabled)
	out = SetBit(out, 3, apu.soundChannels[3].enabled)
	out = SetBit(out, 7, apu.masterEnable)

	return out
}
