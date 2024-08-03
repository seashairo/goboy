package goboy

// @see https://gbdev.io/pandocs/CPU_Registers_and_Flags.html

type CpuRegister byte

const (
	R_NONE CpuRegister = iota

	// The Game Boy has 8 8-bit registers which can be combined to form 4 16-bit
	// registers
	R_A
	// The F reigster is special, and is used to store information about the most
	// recent instruction which has affected flags
	R_F
	R_B
	R_C
	R_D
	R_E
	R_H
	R_L

	// These are the register pairings from 8 to 16-bit registers
	R_AF
	R_BC
	R_DE
	R_HL

	// The stack pointer and program counter are always 16-bit registers and do
	// not have a subdivided part
	R_SP
	R_PC
)

type CpuFlag byte

const (
	// The zero flag is set if and only if the result of an operation is zero.
	// Used by conditional jumps.
	FLAG_Z CpuFlag = 7
	// These flags are used by the DAA instruction only. N indicates whether the
	// previous instruction has been a subtraction, and H indicates carry for the
	// lower 4 bits of the result. DAA also uses the C flag, which must indicate
	// carry for the upper 4 bits. After adding/subtracting two BCD numbers, DAA
	// is used to convert the result to BCD format. BCD numbers range from $00 to
	// $99 rather than $00 to $FF. Because only two flags (C and H) exist to
	// indicate carry-outs of BCD digits, DAA is ineffective for 16-bit operations
	// (which have 4 digits), and use for INC/DEC operations (which do not affect
	// C-flag) has limits.
	FLAG_N CpuFlag = 6
	FLAG_H CpuFlag = 5
	// The carry flag is set in these cases:
	// When the result of an 8-bit addition is higher than $FF.
	// When the result of a 16-bit addition is higher than $FFFF.
	// When the result of a subtraction or comparison is lower than zero (like in
	// 	Z80 and x86 CPUs, but unlike in 65XX and ARM CPUs).
	// When a rotate/shift operation shifts out a “1” bit.
	// Used by conditional jumps and instructions such as ADC, SBC, RL, RLA, etc.
	FLAG_C CpuFlag = 4
)

type CpuRegisters struct {
	a byte
	f byte
	b byte
	c byte
	d byte
	e byte
	h byte
	l byte

	sp uint16
	pc uint16
}

func NewCpuRegisters() *CpuRegisters {
	return &CpuRegisters{
		a:  0x01,
		f:  0xB0,
		b:  0x00,
		c:  0x13,
		d:  0x00,
		e:  0xD8,
		h:  0x01,
		l:  0x4D,
		sp: 0xFFFE,
		pc: 0x0100,
	}
}

func (registers *CpuRegisters) read(register CpuRegister) uint16 {
	switch register {
	case R_A:
		return uint16(registers.a)
	case R_F:
		return uint16(registers.f)
	case R_B:
		return uint16(registers.b)
	case R_C:
		return uint16(registers.c)
	case R_D:
		return uint16(registers.d)
	case R_E:
		return uint16(registers.e)
	case R_H:
		return uint16(registers.h)
	case R_L:
		return uint16(registers.l)
	case R_AF:
		return BytesToUint16(registers.a, registers.f)
	case R_BC:
		return BytesToUint16(registers.b, registers.c)
	case R_DE:
		return BytesToUint16(registers.d, registers.e)
	case R_HL:
		return BytesToUint16(registers.h, registers.l)
	case R_PC:
		return registers.pc
	case R_SP:
		return registers.sp
	default:
		return 0
	}
}

func (registers *CpuRegisters) write(register CpuRegister, value uint16) {
	switch register {
	case R_A:
		registers.a = byte(value)
	case R_F:
		registers.f = byte(value) & 0xF0
	case R_B:
		registers.b = byte(value)
	case R_C:
		registers.c = byte(value)
	case R_D:
		registers.d = byte(value)
	case R_E:
		registers.e = byte(value)
	case R_H:
		registers.h = byte(value)
	case R_L:
		registers.l = byte(value)
	case R_AF:
		hi, lo := Uint16ToBytes(value)
		registers.a = hi
		registers.f = lo & 0xF0
	case R_BC:
		registers.b, registers.c = Uint16ToBytes(value)
	case R_DE:
		registers.d, registers.e = Uint16ToBytes(value)
	case R_HL:
		registers.h, registers.l = Uint16ToBytes(value)
	case R_PC:
		registers.pc = value
	case R_SP:
		registers.sp = value
	default:
		break
	}
}

func (registers *CpuRegisters) setFlag(flag CpuFlag, value bool) {
	registers.f = SetBit(registers.f, byte(flag), value)
}

func (registers *CpuRegisters) readFlag(flag CpuFlag) bool {
	return GetBit(registers.f, byte(flag))
}

func (registers *CpuRegisters) setFlags(z bool, n bool, h bool, c bool) {
	registers.setFlag(FLAG_Z, z)
	registers.setFlag(FLAG_N, n)
	registers.setFlag(FLAG_H, h)
	registers.setFlag(FLAG_C, c)
}

func decodeRegister(b byte) CpuRegister {
	if b > 0x07 {
		return R_NONE
	}

	return [0x08]CpuRegister{
		R_B, R_C, R_D, R_E, R_H, R_L, R_HL, R_A,
	}[b]
}
