package goboy

// Represents a 16-bit register. Each register can be accessed as a single
// 16-bit value or as two 8-bit values. For example, the AF register can be read
// as a single 16-bit value or as two 8-bit values: A and F
type register struct {
	hi byte
	lo byte
}

// Creates a new register from a 16-bit value. By default, registers are
// initialized with 0 values which are then overwritten by the boot rom
func NewRegister(value uint16) register {
	r := register{0, 0}
	r.Set(value)

	return r
}

// Returns the high byte of the register
func (r *register) Hi() byte {
	return r.hi
}

// Returns the low byte of the register
func (r *register) Lo() byte {
	return r.lo
}

// Returns the register as a single 16-bit value
func (r *register) Value() uint16 {
	return BytesToUint16(r.hi, r.lo)
}

// Sets the hi byte of the register
func (r *register) SetHi(value byte) {
	r.hi = value
}

// Sets the lo byte of the register
func (r *register) SetLo(value byte) {
	r.lo = value
}

// Sets the register as a single 16-bit value by setting both the hi and lo
// bytes of the register
func (r *register) Set(value uint16) {
	hi, lo := Uint16ToBytes(value)
	r.SetHi(hi)
	r.SetLo(lo)
}

type Registers struct {
	AF register
	BC register
	DE register
	HL register
	SP register
	PC register
}

func NewRegisters() Registers {
	return Registers{
		AF: NewRegister(0x0000),
		BC: NewRegister(0x0000),
		DE: NewRegister(0x0000),
		HL: NewRegister(0x0000),
		SP: NewRegister(0x0000),
		PC: NewRegister(0x0000),
	}
}

// The lower 8 bits of the AF register form the Flags register. The Flags
// register contains information about the result of the most recent instruction
// that has affected flags (e.g. ADD, SUB, INC, DEC, etc.)
type flag byte

const (
	// The zero flag is set if and only if the result of an operation is zero.
	// Used by conditional jumps.
	FLAG_Z flag = 1 << 7
	// These flags are used by the DAA instruction only. N indicates whether the
	// previous instruction has been a subtraction, and H indicates carry for the
	// lower 4 bits of the result. DAA also uses the C flag, which must indicate
	// carry for the upper 4 bits. After adding/subtracting two BCD numbers, DAA
	// is used to convert the result to BCD format. BCD numbers range from $00 to
	// $99 rather than $00 to $FF. Because only two flags (C and H) exist to
	// indicate carry-outs of BCD digits, DAA is ineffective for 16-bit operations
	// (which have 4 digits), and use for INC/DEC operations (which do not affect
	// C-flag) has limits.
	FLAG_N flag = 1 << 6
	FLAG_H flag = 1 << 5
	// The carry flag is set in these cases:
	// When the result of an 8-bit addition is higher than $FF.
	// When the result of a 16-bit addition is higher than $FFFF.
	// When the result of a subtraction or comparison is lower than zero (like in
	// 	Z80 and x86 CPUs, but unlike in 65XX and ARM CPUs).
	// When a rotate/shift operation shifts out a “1” bit.
	// Used by conditional jumps and instructions such as ADC, SBC, RL, RLA, etc.
	FLAG_C flag = 1 << 4
)

func (r *Registers) SetFlag(flag flag, value bool) {
	if value {
		r.AF.lo |= byte(flag)
	} else {
		r.AF.lo &= byte(^flag)
	}
}

func (r *Registers) GetFlag(flag flag) bool {
	return r.AF.lo&byte(flag) == 1
}

func (r *Registers) ResetFlags() {
	r.AF.SetLo(0)
}
