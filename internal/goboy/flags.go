package goboy

// @see https://gbdev.io/pandocs/CPU_Registers_and_Flags.html
type Flags struct {
	Z bool // zero
	N bool // subtraction
	H bool // half carry
	C bool // carry
}

func NewFlags() Flags {
	return Flags{
		Z: false,
		N: false,
		H: false,
		C: false,
	}
}

func (flags *Flags) AsByte() byte {
	return (b2b(flags.C) << 3) |
		(b2b(flags.H) << 2) |
		(b2b(flags.N) << 1) |
		b2b(flags.Z)
}

func b2b(b bool) byte {
	if b {
		return 1
	}

	return 0
}
