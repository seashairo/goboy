package goboy

type VolumeEnvelope struct {
	finished      bool
	timer         byte
	initialVolume byte
	addMode       bool
	period        byte
	volume        byte
}

func NewVolumeEnvelope() *VolumeEnvelope {
	ve := &VolumeEnvelope{}
	ve.PowerOff()
	return ve
}

func (ve *VolumeEnvelope) Tick() {
	if ve.finished {
		return
	}

	ve.timer--
	if ve.timer <= 0 {
		if ve.period == 0 {
			ve.timer = 8
		} else {
			ve.timer = ve.period
		}

		if ve.addMode && ve.volume < 15 {
			ve.volume++
		} else if !ve.addMode && ve.volume > 0 {
			ve.volume--
		}

		if ve.volume == 0 || ve.volume == 15 {
			ve.finished = true
		}
	}
}

func (ve *VolumeEnvelope) PowerOff() {
	ve.finished = true
	ve.timer = 0
	ve.initialVolume = 0
	ve.addMode = false
	ve.period = 0
	ve.volume = 0
}

func (ve *VolumeEnvelope) SetNR2(value byte) {
	ve.initialVolume = value >> 4
	ve.addMode = GetBit(value, 3)
	ve.period = value & 0b111
}

func (ve *VolumeEnvelope) GetNR2() byte {
	out := (ve.initialVolume << 4) | ve.period
	return SetBit(out, 3, ve.addMode)
}

func (ve *VolumeEnvelope) GetVolume() byte {
	if ve.period > 0 {
		return ve.volume
	}

	return ve.initialVolume
}

func (ve *VolumeEnvelope) Trigger() {
	ve.volume = ve.initialVolume
	ve.finished = false

	if ve.period == 0 {
		ve.timer = 8
	} else {
		ve.timer = ve.period
	}
}
