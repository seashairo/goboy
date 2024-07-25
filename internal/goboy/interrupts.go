package goboy

type InterruptKind byte

const (
	INT_VBLANK = 1 << iota
	INT_LCD
	INT_TIMER
	INT_SERIAL
	INT_JOYPAD
)

type InterruptRegister struct {
	data byte
}

func NewInterruptRegister(data byte) InterruptRegister {
	return InterruptRegister{data: data}
}

func (ie *InterruptRegister) CheckInterrupt(kind InterruptKind) bool {
	return (ie.data & byte(kind)) == 1
}

func (ie *InterruptRegister) SetInterrupt(kind InterruptKind, on bool) {
	ie.data = ie.data | byte(kind)
}

func (ie *InterruptRegister) readByte() byte {
	return ie.data
}

func (ie *InterruptRegister) writeByte(value byte) {
	ie.data = value
}

func (cpu *CPU) handleInterrupts() {
	interruptKinds := [5]InterruptKind{
		INT_VBLANK,
		INT_LCD,
		INT_TIMER,
		INT_SERIAL,
		INT_JOYPAD,
	}

	for _, kind := range interruptKinds {
		if cpu.handleInterrupt(kind) {
			return
		}
	}
}

func (cpu *CPU) checkInterrupt(kind InterruptKind) bool {
	interruptFlags := cpu.bus.io.interrupts.readByte()
	interruptEnabled := cpu.bus.interruptEnableRegister.readByte()

	return (interruptFlags&byte(kind))&(interruptEnabled&byte(kind)) == 1
}

func (cpu *CPU) handleInterrupt(kind InterruptKind) bool {
	if !cpu.checkInterrupt(kind) {
		return false
	}

	push(cpu, R_PC)
	cpu.bus.io.interrupts.SetInterrupt(kind, false)
	cpu.interruptMasterEnabled = false

	return true
}
