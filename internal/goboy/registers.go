package goboy

type Registers struct {
	A              byte
	F              Flags
	B              byte
	C              byte
	D              byte
	E              byte
	H              byte
	L              byte
	ProgramCounter uint16
	StackPointer   uint16
}

func NewRegisters() Registers {
	return Registers{
		A:              0,
		F:              NewFlags(),
		B:              0,
		C:              0,
		D:              0,
		E:              0,
		H:              0,
		L:              0,
		ProgramCounter: 0,
		StackPointer:   0,
	}
}

func (registers Registers) AF() uint16 {
	return BytesToUint16(registers.A, registers.F.AsByte())
}

func (registers Registers) BC() uint16 {
	return BytesToUint16(registers.B, registers.C)
}

func (registers Registers) DE() uint16 {
	return BytesToUint16(registers.D, registers.E)
}

func (registers Registers) HL() uint16 {
	return BytesToUint16(registers.H, registers.L)
}
