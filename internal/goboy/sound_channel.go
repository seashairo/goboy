package goboy

type envDir bool

var (
	envUp   = envDir(true)
	envDown = envDir(false)
)

const (
	squareSoundType = 0
	waveSoundType   = 1
	noiseSoundType  = 2
)

type SoundChannel struct {
	soundType uint8

	enabled        bool
	rightSpeakerOn bool
	leftSpeakerOn  bool

	envelopeDirection   envDir
	envelopeStartVolume byte
	envelopeSweepPace   byte
	envelopeVolume      byte
	envelopeCounter     byte

	t                uint32
	frequencyDivider uint32
	period           uint16

	sweepCounter   byte
	sweepDirection bool
	sweepTime      byte
	sweepShift     byte

	lengthData    byte
	currentLength byte

	waveDuty           byte
	waveDutySeqCounter byte

	waveOutLvl        byte
	wavePatternRAM    [16]byte
	wavePatternCursor byte

	polyFeedbackReg  uint16
	polyDivisorShift byte
	polyDivisorBase  byte
	poly7BitMode     bool
	polySample       byte

	playsContinuously bool
	restartRequested  bool
}

func (sc *SoundChannel) TickFrequency() {
	sc.t += 2 // currently called at 2MHz, so tick twice

	if sc.t >= sc.frequencyDivider {
		sc.t = 0
		switch sc.soundType {
		case squareSoundType:
			sc.waveDutySeqCounter = (sc.waveDutySeqCounter + 1) & 7
		case waveSoundType:
			sc.wavePatternCursor = (sc.wavePatternCursor + 1) & 31
		case noiseSoundType:
			sc.updatePolyCounter()
		}
	}
}

func (sc *SoundChannel) updatePolyCounter() {
	newHigh := (sc.polyFeedbackReg & 0x01) ^ ((sc.polyFeedbackReg >> 1) & 0x01)
	sc.polyFeedbackReg >>= 1
	sc.polyFeedbackReg &^= 1 << 14
	sc.polyFeedbackReg |= newHigh << 14

	if sc.poly7BitMode {
		sc.polyFeedbackReg &^= 1 << 6
		sc.polyFeedbackReg |= newHigh << 6
	}

	if sc.polyFeedbackReg&0x01 == 0 {
		sc.polySample = 1
	} else {
		sc.polySample = 0
	}
}

func (sc *SoundChannel) TickLength() {
	if sc.currentLength > 0 && !sc.playsContinuously {
		sc.currentLength--
		if sc.currentLength == 0 {
			sc.enabled = false
		}
	}

	if sc.restartRequested {
		sc.enabled = true
		sc.restartRequested = false

		if sc.lengthData == 0 {
			if sc.soundType == waveSoundType {
				sc.lengthData = 255
			} else {
				sc.lengthData = 64
			}
		}

		sc.currentLength = sc.lengthData
		sc.envelopeVolume = sc.envelopeStartVolume
		sc.sweepCounter = 0
		sc.wavePatternCursor = 0
		sc.polyFeedbackReg = 0xFFFF
	}
}

func (sc *SoundChannel) TickSweep() {
	if sc.sweepTime != 0 {
		if sc.sweepCounter < sc.sweepTime {
			sc.sweepCounter++
		} else {
			sc.sweepCounter = 0
			var nextFreq uint16

			if sc.sweepDirection {
				nextFreq = sc.period - (sc.period >> uint16(sc.sweepShift))
			} else {
				nextFreq = sc.period + (sc.period >> uint16(sc.sweepShift))
			}

			if nextFreq > 2047 {
				sc.enabled = false
			} else {
				sc.period = nextFreq
				sc.updateFrequency()
			}
		}
	}
}

func (sc *SoundChannel) TickVolumeEnvelope() {
	if sc.envelopeSweepPace != 0 {
		if sc.envelopeCounter < sc.envelopeSweepPace {
			sc.envelopeCounter++
		} else {
			sc.envelopeCounter = 0

			if sc.envelopeDirection == envUp && sc.envelopeVolume < 15 {
				sc.envelopeVolume++

			} else if sc.envelopeDirection == envDown && sc.envelopeVolume > 0 {
				sc.envelopeVolume--
			}
		}
	}
}

var dutyCycleTable = [4][8]byte{
	{0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 1, 1, 1},
	{0, 1, 1, 1, 1, 1, 1, 0},
}

func (sc *SoundChannel) inDutyCycle() bool {
	sel := sc.waveDuty
	counter := sc.waveDutySeqCounter
	return dutyCycleTable[sel][counter] == 1
}

func (sc *SoundChannel) getSample() (byte, byte) {
	sample := byte(0)
	if sc.enabled {
		switch sc.soundType {
		case squareSoundType:
			if sc.inDutyCycle() {
				sample = sc.envelopeVolume
			} else {
				sample = 0
			}

		case waveSoundType:
			if sc.waveOutLvl > 0 {
				sampleByte := sc.wavePatternRAM[sc.wavePatternCursor/2]
				if sc.wavePatternCursor&1 == 0 {
					sample = sampleByte >> 4
				} else {
					sample = sampleByte & 0x0f
				}
			}

		case noiseSoundType:
			if sc.frequencyDivider > 0 {
				sample = sc.envelopeVolume * sc.polySample
			}
		}
	}

	left, right := byte(0), byte(0)

	if sc.leftSpeakerOn {
		left = sample
	}

	if sc.rightSpeakerOn {
		right = sample
	}

	return left, right
}

func (sc *SoundChannel) updateFrequency() {
	switch sc.soundType {
	case waveSoundType:
		sc.frequencyDivider = 2 * (2048 - uint32(sc.period))
	case noiseSoundType:
		divider := uint32(8)

		if sc.polyDivisorBase > 0 {
			if sc.polyDivisorShift < 14 {
				divider = uint32(sc.polyDivisorBase) << uint32(sc.polyDivisorShift+4)
			} else {
				divider = 0
			}
		}

		sc.frequencyDivider = divider
	case squareSoundType:
		sc.frequencyDivider = 4 * (2048 - uint32(sc.period))
	}
}

func (sc *SoundChannel) writeWaveOnOffReg(value byte) {
	sc.enabled = GetBit(value, 7)
}

func (sc *SoundChannel) writeWavePatternValue(addr uint16, value byte) {
	sc.wavePatternRAM[addr] = value
}

func (sc *SoundChannel) writePolyCounterReg(value byte) {
	sc.poly7BitMode = GetBit(value, 3)
	sc.polyDivisorShift = value >> 4
	sc.polyDivisorBase = value & 0x07
}

func (sc *SoundChannel) readPolyCounterReg() byte {
	value := SetBit(0, 3, sc.poly7BitMode)
	value |= sc.polyDivisorShift << 4
	value |= sc.polyDivisorBase

	return value
}

func (sc *SoundChannel) writeWaveOutLvlReg(value byte) {
	sc.waveOutLvl = (value >> 5) & 0x03
}

func (sc *SoundChannel) readWaveOnOffReg() byte {
	return SetBit(0b01111111, 7, sc.enabled)
}

func (sc *SoundChannel) readWaveOutLvlReg() byte {
	return (sc.waveOutLvl << 5) | 0x9f
}

func (sc *SoundChannel) writeLengthDataReg(value byte) {
	switch sc.soundType {
	case waveSoundType:
		sc.lengthData = 255 - value
	case noiseSoundType:
		sc.lengthData = 64 - (value & 0x3f)
	default:
		panic("writeLengthData: unexpected sound type")
	}
}

func (sc *SoundChannel) readLengthDataReg() byte {
	switch sc.soundType {
	case waveSoundType:
		return 255 - sc.lengthData
	case noiseSoundType:
		return 64 - sc.lengthData
	default:
		panic("readLengthData: unexpected sound type")
	}
}

func (sc *SoundChannel) writeLenDutyReg(value byte) {
	sc.lengthData = 64 - value&0x3f
	sc.waveDuty = value >> 6
}

func (sc *SoundChannel) readLenDutyReg() byte {
	return (sc.waveDuty << 6) | 0x3f
}

func (sc *SoundChannel) writeSweepReg(value byte) {
	sc.sweepTime = (value >> 4) & 0x07
	sc.sweepShift = value & 0x07
	sc.sweepDirection = GetBit(value, 3)
}

func (sc *SoundChannel) readSweepReg() byte {
	value := sc.sweepTime << 4
	value |= sc.sweepShift
	value = SetBit(value, 3, bool(sc.sweepDirection))
	return SetBit(value, 7, true)
}

func (sc *SoundChannel) writeVolumeEnvelope(value byte) {
	sc.envelopeStartVolume = value >> 4

	if sc.envelopeStartVolume == 0 {
		sc.enabled = false
	}

	sc.envelopeDirection = envDir(GetBit(value, 3))

	sc.envelopeSweepPace = value & 0x07
}

func (sc *SoundChannel) readVolumeEnvelope() byte {
	value := sc.envelopeStartVolume<<4 | sc.envelopeSweepPace
	return SetBit(value, 3, bool(sc.envelopeDirection))
}

func (sc *SoundChannel) writePeriodLow(value byte) {
	sc.period &^= 0x00ff
	sc.period |= uint16(value)
	sc.updateFrequency()
}

func (sc *SoundChannel) readPeriodLow() byte {
	return 0xFF
}

func (sc *SoundChannel) writePeriodHigh(value byte) {
	if value&0x80 != 0 {
		sc.restartRequested = true
	}

	sc.playsContinuously = value&0x40 == 0
	sc.period &^= 0xFF00
	sc.period |= uint16(value&0x07) << 8
	sc.updateFrequency()
}

func (sc *SoundChannel) readPeriodHigh() byte {
	value := byte(0xFF)

	if sc.playsContinuously {
		value &^= 0x40
	}

	return value
}
