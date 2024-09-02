package goboy

type InterruptKind byte

const (
	INT_VBLANK = iota
	INT_LCD
	INT_TIMER
	INT_SERIAL
	INT_JOYPAD
)

type InterruptRegister struct {
	data byte
}

func NewInterruptRegister(data byte) *InterruptRegister {
	return &InterruptRegister{data: data}
}

func (ie *InterruptRegister) SetInterrupt(kind InterruptKind, on bool) {
	ie.data = SetBit(ie.data, byte(kind), on)
}

func (ie *InterruptRegister) readByte() byte {
	return ie.data
}

func (ie *InterruptRegister) writeByte(value byte) {
	ie.data = value
}

func (cpu *CPU) handleInterrupts() {
	if cpu.handleInterrupt(INT_VBLANK, 0x40) {
		return
	}

	if cpu.handleInterrupt(INT_LCD, 0x48) {
		return
	}

	if cpu.handleInterrupt(INT_TIMER, 0x50) {
		return
	}

	if cpu.handleInterrupt(INT_SERIAL, 0x58) {
		return
	}

	if cpu.handleInterrupt(INT_JOYPAD, 0x60) {
		return
	}
}

func (cpu *CPU) checkInterrupt(kind InterruptKind) bool {
	interruptFlags := cpu.bus.readByte(IO_IF)
	ieRegister := cpu.bus.readByte(INTERRUPT_ENABLE_REGISTER_START)

	interruptFlagged := GetBit(interruptFlags, byte(kind))
	interruptEnabled := GetBit(ieRegister, byte(kind))

	return interruptFlagged && interruptEnabled
}

func (cpu *CPU) handleInterrupt(kind InterruptKind, address uint16) bool {
	if !cpu.checkInterrupt(kind) {
		return false
	}

	cpu.gameboy.ClearInterrupt(kind)
	cpu.interruptMasterEnabled = false
	push(cpu, R_PC)
	cpu.registers.write(R_PC, address)

	return true
}
