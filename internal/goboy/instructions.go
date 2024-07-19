package goboy

import "fmt"

type condition byte

const (
	C_ANY condition = iota
	C_Z
	C_NZ
	C_C
	C_NC
)

func checkCondition(cpu *CPU, cond condition) bool {
	switch cond {
	case C_ANY:
		return true
	case C_Z:
		return cpu.registers.readFlag(FLAG_Z)
	case C_NZ:
		return !cpu.registers.readFlag(FLAG_Z)
	case C_C:
		return cpu.registers.readFlag(FLAG_C)
	case C_NC:
		return !cpu.registers.readFlag(FLAG_C)
	}

	panic("Encountered unknown condition")
}

type instruction func(cpu *CPU)

func fetchInstruction(opcode byte) instruction {
	instruction := instructions[opcode]

	if instruction != nil {
		return instruction
	}

	panic(fmt.Sprintf("No instruction found for opcode 0x%2.2X", opcode))
}

var instructions = [0x100]instruction{
	0x00: func(_ *CPU) {},
	0x01: func(cpu *CPU) {
		ldN16ToR16(cpu, R_BC)
	},
	0x02: func(cpu *CPU) {
		ldR8ToMR16(cpu, R_A, R_BC)
	},
	0x03: func(cpu *CPU) {
		incR16(cpu, R_BC)
	},
	0x04: func(cpu *CPU) {
		incR8(cpu, R_B)
	},
	0x05: func(cpu *CPU) {
		decR8(cpu, R_B)
	},
	0x06: func(cpu *CPU) {
		ldN8ToR8(cpu, R_B)
	},
	0x07: rlca,
	0x08: func(cpu *CPU) {
		ldR16ToA16(cpu, R_SP)
	},
	0x09: func(cpu *CPU) {
		addR16(cpu, R_BC)
	},
	0x0A: func(cpu *CPU) {
		ldMR16ToR8(cpu, R_BC, R_A)
	},
	0x0B: func(cpu *CPU) {
		decR16(cpu, R_BC)
	},
	0x0C: func(cpu *CPU) {
		incR8(cpu, R_C)
	},
	0x0D: func(cpu *CPU) {
		decR8(cpu, R_C)
	},
	0x0E: func(cpu *CPU) {
		ldN8ToR8(cpu, R_C)
	},
	0x0F: rrca,
	0x10: func(cpu *CPU) {
		// do nothing?
		// todo: what does STOP actually do?
	},
	0x11: func(cpu *CPU) {
		ldN16ToR16(cpu, R_DE)
	},
	0x12: func(cpu *CPU) {
		ldR8ToMR16(cpu, R_A, R_DE)
	},
	0x13: func(cpu *CPU) {
		incR16(cpu, R_DE)
	},
	0x14: func(cpu *CPU) {
		incR8(cpu, R_D)
	},
	0x15: func(cpu *CPU) {
		decR8(cpu, R_D)
	},
	0x16: func(cpu *CPU) {
		ldN8ToR8(cpu, R_D)
	},
	0x17: rla,
	0x18: func(cpu *CPU) {
		jr(cpu, C_ANY)
	},
	0x19: func(cpu *CPU) {
		addR16(cpu, R_DE)
	},
	0x1A: func(cpu *CPU) {
		ldMR16ToR8(cpu, R_DE, R_A)
	},
	0x1B: func(cpu *CPU) {
		decR16(cpu, R_DE)
	},
	0x1C: func(cpu *CPU) {
		incR8(cpu, R_E)
	},
	0x1D: func(cpu *CPU) {
		decR8(cpu, R_E)
	},
	0x1E: func(cpu *CPU) {
		ldN8ToR8(cpu, R_E)
	},
	0x1F: rra,
	0x20: func(cpu *CPU) {
		jr(cpu, C_NZ)
	},
	0x21: func(cpu *CPU) {
		ldN16ToR16(cpu, R_HL)
	},
	0x22: func(cpu *CPU) {
		ldR8ToMR16(cpu, R_A, R_HL)
		incR16(cpu, R_HL)
	},
	0x23: func(cpu *CPU) {
		incR16(cpu, R_HL)
	},
	0x24: func(cpu *CPU) {
		incR8(cpu, R_H)
	},
	0x25: func(cpu *CPU) {
		decR8(cpu, R_H)
	},
	0x26: func(cpu *CPU) {
		ldN8ToR8(cpu, R_H)
	},
	0x27: daa,
	0x28: func(cpu *CPU) {
		jr(cpu, C_Z)
	},
	0x29: func(cpu *CPU) {
		addR16(cpu, R_HL)
	},
	0x2A: func(cpu *CPU) {
		ldMR16ToR8(cpu, R_HL, R_A)
		incR16(cpu, R_HL)
	},
	0x2B: func(cpu *CPU) {
		decR16(cpu, R_HL)
	},
	0x2C: func(cpu *CPU) {
		incR8(cpu, R_L)
	},
	0x2D: func(cpu *CPU) {
		decR8(cpu, R_L)
	},
	0x2E: func(cpu *CPU) {
		ldN8ToR8(cpu, R_L)
	},
	0x2F: cpl,
	0x30: func(cpu *CPU) {
		jr(cpu, C_NC)
	},
	0x31: func(cpu *CPU) {
		ldN16ToR16(cpu, R_SP)
	},
	0x32: func(cpu *CPU) {
		ldR8ToMR16(cpu, R_A, R_HL)
		decR16(cpu, R_HL)
	},
	0x33: func(cpu *CPU) {
		incR16(cpu, R_SP)
	},
	0x34: func(cpu *CPU) {
		incMR16(cpu, R_HL)
	},
	0x35: func(cpu *CPU) {
		decMR16(cpu, R_HL)
	},
	0x36: func(cpu *CPU) {
		ldN8ToMR16(cpu, R_HL)
	},
	0x37: scf,
	0x38: func(cpu *CPU) {
		jr(cpu, C_C)
	},
	0x39: func(cpu *CPU) {
		addR16(cpu, R_SP)
	},
	0x3A: func(cpu *CPU) {
		ldMR16ToR8(cpu, R_HL, R_A)
		decR16(cpu, R_HL)
	},
	0x3B: func(cpu *CPU) {
		decR16(cpu, R_SP)
	},
	0x3C: func(cpu *CPU) {
		incR8(cpu, R_A)
	},
	0x3D: func(cpu *CPU) {
		decR8(cpu, R_A)
	},
	0x3E: func(cpu *CPU) {
		ldN8ToR8(cpu, R_A)
	},
	0x3F: ccf,
	0x40: func(cpu *CPU) {
		ldR8ToR8(cpu, R_B, R_B)
	},
	0x41: func(cpu *CPU) {
		ldR8ToR8(cpu, R_C, R_B)
	},
	0x42: func(cpu *CPU) {
		ldR8ToR8(cpu, R_D, R_B)
	},
	0x43: func(cpu *CPU) {
		ldR8ToR8(cpu, R_E, R_B)
	},
	0x44: func(cpu *CPU) {
		ldR8ToR8(cpu, R_H, R_B)
	},
	0x45: func(cpu *CPU) {
		ldR8ToR8(cpu, R_L, R_B)
	},
	0x46: func(cpu *CPU) {
		ldMR16ToR8(cpu, R_HL, R_B)
	},
	0x47: func(cpu *CPU) {
		ldR8ToR8(cpu, R_A, R_B)
	},
	0x48: func(cpu *CPU) {
		ldR8ToR8(cpu, R_B, R_C)
	},
	0x49: func(cpu *CPU) {
		ldR8ToR8(cpu, R_C, R_C)
	},
	0x4A: func(cpu *CPU) {
		ldR8ToR8(cpu, R_D, R_C)
	},
	0x4B: func(cpu *CPU) {
		ldR8ToR8(cpu, R_E, R_C)
	},
	0x4C: func(cpu *CPU) {
		ldR8ToR8(cpu, R_H, R_C)
	},
	0x4D: func(cpu *CPU) {
		ldR8ToR8(cpu, R_L, R_C)
	},
	0x4E: func(cpu *CPU) {
		ldMR16ToR8(cpu, R_HL, R_C)
	},
	0x4F: func(cpu *CPU) {
		ldR8ToR8(cpu, R_A, R_C)
	},
	0x50: func(cpu *CPU) {
		ldR8ToR8(cpu, R_B, R_D)
	},
	0x51: func(cpu *CPU) {
		ldR8ToR8(cpu, R_C, R_D)
	},
	0x52: func(cpu *CPU) {
		ldR8ToR8(cpu, R_D, R_D)
	},
	0x53: func(cpu *CPU) {
		ldR8ToR8(cpu, R_E, R_D)
	},
	0x54: func(cpu *CPU) {
		ldR8ToR8(cpu, R_H, R_D)
	},
	0x55: func(cpu *CPU) {
		ldR8ToR8(cpu, R_L, R_D)
	},
	0x56: func(cpu *CPU) {
		ldMR16ToR8(cpu, R_HL, R_D)
	},
	0x57: func(cpu *CPU) {
		ldR8ToR8(cpu, R_A, R_D)
	},
	0x58: func(cpu *CPU) {
		ldR8ToR8(cpu, R_B, R_E)
	},
	0x59: func(cpu *CPU) {
		ldR8ToR8(cpu, R_C, R_E)
	},
	0x5A: func(cpu *CPU) {
		ldR8ToR8(cpu, R_D, R_E)
	},
	0x5B: func(cpu *CPU) {
		ldR8ToR8(cpu, R_E, R_E)
	},
	0x5C: func(cpu *CPU) {
		ldR8ToR8(cpu, R_H, R_E)
	},
	0x5D: func(cpu *CPU) {
		ldR8ToR8(cpu, R_L, R_E)
	},
	0x5E: func(cpu *CPU) {
		ldMR16ToR8(cpu, R_HL, R_E)
	},
	0x5F: func(cpu *CPU) {
		ldR8ToR8(cpu, R_A, R_E)
	},
	0x60: func(cpu *CPU) {
		ldR8ToR8(cpu, R_B, R_H)
	},
	0x61: func(cpu *CPU) {
		ldR8ToR8(cpu, R_C, R_H)
	},
	0x62: func(cpu *CPU) {
		ldR8ToR8(cpu, R_D, R_H)
	},
	0x63: func(cpu *CPU) {
		ldR8ToR8(cpu, R_E, R_H)
	},
	0x64: func(cpu *CPU) {
		ldR8ToR8(cpu, R_H, R_H)
	},
	0x65: func(cpu *CPU) {
		ldR8ToR8(cpu, R_L, R_H)
	},
	0x66: func(cpu *CPU) {
		ldMR16ToR8(cpu, R_HL, R_H)
	},
	0x67: func(cpu *CPU) {
		ldR8ToR8(cpu, R_A, R_H)
	},
	0x68: func(cpu *CPU) {
		ldR8ToR8(cpu, R_B, R_L)
	},
	0x69: func(cpu *CPU) {
		ldR8ToR8(cpu, R_C, R_L)
	},
	0x6A: func(cpu *CPU) {
		ldR8ToR8(cpu, R_D, R_L)
	},
	0x6B: func(cpu *CPU) {
		ldR8ToR8(cpu, R_E, R_L)
	},
	0x6C: func(cpu *CPU) {
		ldR8ToR8(cpu, R_H, R_L)
	},
	0x6D: func(cpu *CPU) {
		ldR8ToR8(cpu, R_L, R_L)
	},
	0x6E: func(cpu *CPU) {
		ldMR16ToR8(cpu, R_HL, R_L)
	},
	0x6F: func(cpu *CPU) {
		ldR8ToR8(cpu, R_A, R_L)
	},
	0x70: func(cpu *CPU) {
		ldR8ToMR16(cpu, R_B, R_HL)
	},
	0x71: func(cpu *CPU) {
		ldR8ToMR16(cpu, R_C, R_HL)
	},
	0x72: func(cpu *CPU) {
		ldR8ToMR16(cpu, R_D, R_HL)
	},
	0x73: func(cpu *CPU) {
		ldR8ToMR16(cpu, R_E, R_HL)
	},
	0x74: func(cpu *CPU) {
		ldR8ToMR16(cpu, R_H, R_HL)
	},
	0x75: func(cpu *CPU) {
		ldR8ToMR16(cpu, R_L, R_HL)
	},
	0x76: func(cpu *CPU) {
		// todo: HALT
	},
	0x77: func(cpu *CPU) {
		ldR8ToMR16(cpu, R_A, R_HL)
	},
	0x78: func(cpu *CPU) {
		ldR8ToR8(cpu, R_B, R_A)
	},
	0x79: func(cpu *CPU) {
		ldR8ToR8(cpu, R_C, R_A)
	},
	0x7A: func(cpu *CPU) {
		ldR8ToR8(cpu, R_D, R_A)
	},
	0x7B: func(cpu *CPU) {
		ldR8ToR8(cpu, R_E, R_A)
	},
	0x7C: func(cpu *CPU) {
		ldR8ToR8(cpu, R_H, R_A)
	},
	0x7D: func(cpu *CPU) {
		ldR8ToR8(cpu, R_L, R_A)
	},
	0x7E: func(cpu *CPU) {
		ldMR16ToR8(cpu, R_HL, R_A)
	},
	0x7F: func(cpu *CPU) {
		ldR8ToR8(cpu, R_A, R_A)
	},
	0x80: func(cpu *CPU) {
		addR8(cpu, R_B)
	},
	0x81: func(cpu *CPU) {
		addR8(cpu, R_C)
	},
	0x82: func(cpu *CPU) {
		addR8(cpu, R_D)
	},
	0x83: func(cpu *CPU) {
		addR8(cpu, R_E)
	},
	0x84: func(cpu *CPU) {
		addR8(cpu, R_H)
	},
	0x85: func(cpu *CPU) {
		addR8(cpu, R_L)
	},
	0x86: func(cpu *CPU) {
		addMR16(cpu, R_HL)
	},
	0x87: func(cpu *CPU) {
		addR8(cpu, R_A)
	},
	0x88: func(cpu *CPU) {
		adcR8(cpu, R_B)
	},
	0x89: func(cpu *CPU) {
		adcR8(cpu, R_C)
	},
	0x8A: func(cpu *CPU) {
		adcR8(cpu, R_D)
	},
	0x8B: func(cpu *CPU) {
		adcR8(cpu, R_E)
	},
	0x8C: func(cpu *CPU) {
		adcR8(cpu, R_H)
	},
	0x8D: func(cpu *CPU) {
		adcR8(cpu, R_L)
	},
	0x8E: func(cpu *CPU) {
		adcMR16(cpu, R_HL)
	},
	0x8F: func(cpu *CPU) {
		adcR8(cpu, R_A)
	},
	0x90: func(cpu *CPU) {
		subR8(cpu, R_B)
	},
	0x91: func(cpu *CPU) {
		subR8(cpu, R_C)
	},
	0x92: func(cpu *CPU) {
		subR8(cpu, R_D)
	},
	0x93: func(cpu *CPU) {
		subR8(cpu, R_E)
	},
	0x94: func(cpu *CPU) {
		subR8(cpu, R_H)
	},
	0x95: func(cpu *CPU) {
		subR8(cpu, R_L)
	},
	0x96: func(cpu *CPU) {
		subMR16(cpu, R_HL)
	},
	0x97: func(cpu *CPU) {
		subR8(cpu, R_A)
	},
	0x98: func(cpu *CPU) {
		sbcR8(cpu, R_B)
	},
	0x99: func(cpu *CPU) {
		sbcR8(cpu, R_C)
	},
	0x9A: func(cpu *CPU) {
		sbcR8(cpu, R_D)
	},
	0x9B: func(cpu *CPU) {
		sbcR8(cpu, R_E)
	},
	0x9C: func(cpu *CPU) {
		sbcR8(cpu, R_H)
	},
	0x9D: func(cpu *CPU) {
		sbcR8(cpu, R_L)
	},
	0x9E: func(cpu *CPU) {
		sbcMR16(cpu, R_HL)
	},
	0x9F: func(cpu *CPU) {
		sbcR8(cpu, R_A)
	},
	0xA0: func(cpu *CPU) {
		andR8(cpu, R_B)
	},
	0xA1: func(cpu *CPU) {
		andR8(cpu, R_C)
	},
	0xA2: func(cpu *CPU) {
		andR8(cpu, R_D)
	},
	0xA3: func(cpu *CPU) {
		andR8(cpu, R_E)
	},
	0xA4: func(cpu *CPU) {
		andR8(cpu, R_H)
	},
	0xA5: func(cpu *CPU) {
		andR8(cpu, R_L)
	},
	0xA6: func(cpu *CPU) {
		andMR16(cpu, R_HL)
	},
	0xA7: func(cpu *CPU) {
		andR8(cpu, R_A)
	},
	0xA8: func(cpu *CPU) {
		xorR8(cpu, R_B)
	},
	0xA9: func(cpu *CPU) {
		xorR8(cpu, R_C)
	},
	0xAA: func(cpu *CPU) {
		xorR8(cpu, R_D)
	},
	0xAB: func(cpu *CPU) {
		xorR8(cpu, R_E)
	},
	0xAC: func(cpu *CPU) {
		xorR8(cpu, R_H)
	},
	0xAD: func(cpu *CPU) {
		xorR8(cpu, R_L)
	},
	0xAE: func(cpu *CPU) {
		xorMR16(cpu, R_HL)
	},
	0xAF: func(cpu *CPU) {
		xorR8(cpu, R_A)
	},
	0xB0: func(cpu *CPU) {
		orR8(cpu, R_B)
	},
	0xB1: func(cpu *CPU) {
		orR8(cpu, R_C)
	},
	0xB2: func(cpu *CPU) {
		orR8(cpu, R_D)
	},
	0xB3: func(cpu *CPU) {
		orR8(cpu, R_E)
	},
	0xB4: func(cpu *CPU) {
		orR8(cpu, R_H)
	},
	0xB5: func(cpu *CPU) {
		orR8(cpu, R_L)
	},
	0xB6: func(cpu *CPU) {
		orMR16(cpu, R_HL)
	},
	0xB7: func(cpu *CPU) {
		orR8(cpu, R_A)
	},
	0xB8: func(cpu *CPU) {
		cpR8(cpu, R_B)
	},
	0xB9: func(cpu *CPU) {
		cpR8(cpu, R_C)
	},
	0xBA: func(cpu *CPU) {
		cpR8(cpu, R_D)
	},
	0xBB: func(cpu *CPU) {
		cpR8(cpu, R_E)
	},
	0xBC: func(cpu *CPU) {
		cpR8(cpu, R_H)
	},
	0xBD: func(cpu *CPU) {
		cpR8(cpu, R_L)
	},
	0xBE: func(cpu *CPU) {
		cpMR8(cpu, R_HL)
	},
	0xBF: func(cpu *CPU) {
		cpR8(cpu, R_A)
	},
	0xC0: func(cpu *CPU) {
		ret(cpu, C_NZ)
	},
	0xC1: func(cpu *CPU) {
		pop(cpu, R_B, R_C)
	},
	0xC2: func(cpu *CPU) {
		jpA16(cpu, C_NZ)
	},
	0xC3: func(cpu *CPU) {
		nextAddress := readWordFromPC(cpu)
		cpu.registers.write(R_PC, nextAddress)
	},
	0xC4: func(cpu *CPU) {
		call(cpu, C_NZ)
	},
	0xC5: func(cpu *CPU) {
		push(cpu, R_BC)
	},
	0xC6: addN8A,
	0xC7: func(cpu *CPU) {
		rst(cpu, 0x00)
	},
	0xC8: func(cpu *CPU) {
		ret(cpu, C_Z)
	},
	0xC9: func(cpu *CPU) {
		ret(cpu, C_ANY)
	},
	0xCA: func(cpu *CPU) {
		jpA16(cpu, C_Z)
	},
	0xCC: func(cpu *CPU) {
		call(cpu, C_Z)
	},
	0xCD: func(cpu *CPU) {
		call(cpu, C_ANY)
	},
	0xCE: adcN8,
	0xCF: func(cpu *CPU) {
		rst(cpu, 0x08)
	},
	0xD0: func(cpu *CPU) {
		ret(cpu, C_NC)
	},
	0xD1: func(cpu *CPU) {
		pop(cpu, R_D, R_E)
	},
	0xD2: func(cpu *CPU) {
		jpA16(cpu, C_NC)
	},
	0xD3: invalidInstruction,
	0xD4: func(cpu *CPU) {
		call(cpu, C_NC)
	},
	0xD5: func(cpu *CPU) {
		push(cpu, R_DE)
	},
	0xD6: subN8,
	0xD7: func(cpu *CPU) {
		rst(cpu, 0x10)
	},
	0xD8: func(cpu *CPU) {
		ret(cpu, C_C)
	},
	0xD9: func(cpu *CPU) {
		cpu.masterInterruptEnabled = true
		ret(cpu, C_NZ)
	},
	0xDB: invalidInstruction,
	0xDA: func(cpu *CPU) {
		jpA16(cpu, C_C)
	},
	0xDC: func(cpu *CPU) {
		call(cpu, C_NC)
	},
	0xDD: invalidInstruction,
	0xDE: sbcN8,
	0xDF: func(cpu *CPU) {
		rst(cpu, 0x18)
	},
	0xE0: func(cpu *CPU) {
		ldhR8ToA8(cpu, R_A)
	},
	0xE1: func(cpu *CPU) {
		pop(cpu, R_H, R_L)
	},
	0xE2: func(cpu *CPU) {
		ldR8ToMR8(cpu, R_A, R_C)
	},
	0xE3: invalidInstruction,
	0xE4: invalidInstruction,
	0xE5: func(cpu *CPU) {
		push(cpu, R_HL)
	},
	0xE6: andN8,
	0xE7: func(cpu *CPU) {
		rst(cpu, 0x20)
	},
	0xE8: addN8SP,
	0xE9: func(cpu *CPU) {
		jpR16(cpu, R_HL)
	},
	0xEA: func(cpu *CPU) {
		ldR8ToN16(cpu, R_A)
	},
	0xEB: invalidInstruction,
	0xEC: invalidInstruction,
	0xED: invalidInstruction,
	0xEE: xorN8,
	0xEF: func(cpu *CPU) {
		rst(cpu, 0x28)
	},
	0xF0: func(cpu *CPU) {
		ldhA8ToR8(cpu, R_A)
	},
	0xF1: func(cpu *CPU) {
		pop(cpu, R_A, R_F)
	},
	0xF2: func(cpu *CPU) {
		ldMR8ToR8(cpu, R_C, R_A)
	},
	0xF3: func(cpu *CPU) {
		cpu.masterInterruptEnabled = false
	},
	0xF4: invalidInstruction,
	0xF5: func(cpu *CPU) {
		push(cpu, R_AF)
	},
	0xF6: orN8,
	0xF7: func(cpu *CPU) {
		rst(cpu, 0x30)
	},
	0xF8: func(cpu *CPU) {
		ldR16E8ToR16(cpu, R_SP, R_HL)
	},
	0xF9: func(cpu *CPU) {
		ldR16ToR16(cpu, R_SP, R_HL)
	},
	0xFA: func(cpu *CPU) {
		ldA16ToR8(cpu, R_A)
	},
	0xFB: func(cpu *CPU) {
		cpu.masterInterruptEnabled = true
	},
	0xFC: invalidInstruction,
	0xFD: invalidInstruction,
	0xFE: cpN8,
	0xFF: func(cpu *CPU) {
		rst(cpu, 0x38)
	},
}

func readByteFromPC(cpu *CPU) byte {
	return cpu.fetchNextOpcode()
}

func readWordFromPC(cpu *CPU) uint16 {
	lo := readByteFromPC(cpu)
	hi := readByteFromPC(cpu)

	return BytesToUint16(hi, lo)
}

// @see https://gist.github.com/SonoSooS/c0055300670d678b5ae8433e20bea595#opcode-holes-not-implemented-opcodes
func invalidInstruction(cpu *CPU) {
	panic("this crashes the game boy")
}

func ldR8ToR8(cpu *CPU, src CpuRegister, dest CpuRegister) {
	cpu.registers.write(dest, cpu.registers.read(src))
}

func ldR16ToR16(cpu *CPU, src CpuRegister, dest CpuRegister) {
	cpu.registers.write(dest, cpu.registers.read(src))
}

func ldR16E8ToR16(cpu *CPU, src CpuRegister, dest CpuRegister) {
	e8 := readByteFromPC(cpu)
	offset := uint16(int8(e8))

	cpu.registers.write(dest, cpu.registers.read(src)+offset)
}

func ldR8ToMR16(cpu *CPU, src CpuRegister, dest CpuRegister) {
	cpu.bus.writeByte(cpu.registers.read(dest), byte(cpu.registers.read(src)))
}

func ldR16ToA16(cpu *CPU, src CpuRegister) {
	value := cpu.registers.read(src)
	address := readWordFromPC(cpu)

	cpu.bus.writeByte(address, byte(value&0xFF))
	cpu.bus.writeByte(address+1, byte(value>>8))
}

func ldR8ToMR8(cpu *CPU, src CpuRegister, dest CpuRegister) {
	cpu.bus.writeByte(0xFF00+cpu.registers.read(dest), byte(cpu.registers.read(src)))
}

func ldMR16ToR8(cpu *CPU, src CpuRegister, dest CpuRegister) {
	address := cpu.registers.read(src)
	value := cpu.bus.readByte(address)
	cpu.registers.write(dest, uint16(value))
}

func ldMR8ToR8(cpu *CPU, src CpuRegister, dest CpuRegister) {
	address := cpu.registers.read(src)
	value := 0xFF00 + uint16(cpu.bus.readByte(address))
	cpu.registers.write(dest, value)
}

func ldN16ToR16(cpu *CPU, dest CpuRegister) {
	n16 := readWordFromPC(cpu)
	cpu.registers.write(dest, n16)
}

func ldN8ToR8(cpu *CPU, dest CpuRegister) {
	n8 := readByteFromPC(cpu)
	cpu.registers.write(dest, uint16(n8))
}

func ldN8ToMR16(cpu *CPU, dest CpuRegister) {
	n8 := readByteFromPC(cpu)
	address := cpu.registers.read(dest)
	cpu.bus.writeByte(address, n8)
}

func ldR8ToN16(cpu *CPU, src CpuRegister) {
	dest := readWordFromPC(cpu)
	cpu.bus.writeByte(dest, byte(cpu.registers.read(src)))
}

func ldA16ToR8(cpu *CPU, dest CpuRegister) {
	a16 := readWordFromPC(cpu)
	value := cpu.bus.readByte(a16)
	cpu.registers.write(dest, uint16(value))
}

func xorR8(cpu *CPU, src CpuRegister) {
	xor(cpu, byte(cpu.registers.read(src)))
}

func xorMR16(cpu *CPU, src CpuRegister) {
	xor(cpu, cpu.bus.readByte(cpu.registers.read(src)))
}

func xorN8(cpu *CPU) {
	xor(cpu, readByteFromPC(cpu))
}

func xor(cpu *CPU, comparator byte) {
	a := byte(cpu.registers.read(R_A))
	result := comparator ^ a

	cpu.registers.write(R_A, uint16(result))
	cpu.registers.setFlags(result == 0, false, false, false)
}

func incR8(cpu *CPU, reg CpuRegister) {
	value := cpu.registers.read(reg) + 1
	cpu.registers.write(reg, value)

	cpu.registers.setFlag(FLAG_Z, value == 0)
	cpu.registers.setFlag(FLAG_N, false)
	cpu.registers.setFlag(FLAG_H, (value&0x0F) == 0x0F)
}

func incR16(cpu *CPU, reg CpuRegister) {
	cpu.registers.write(reg, cpu.registers.read(reg)+1)
}

func incMR16(cpu *CPU, reg CpuRegister) {
	address := cpu.registers.read(reg)
	value := cpu.bus.readByte(address) + 1
	cpu.bus.writeByte(address, value)

	cpu.registers.setFlag(FLAG_Z, value == 0)
	cpu.registers.setFlag(FLAG_N, false)
	cpu.registers.setFlag(FLAG_H, (value&0x0F) == 0x0F)
}

func decR8(cpu *CPU, reg CpuRegister) {
	value := cpu.registers.read(reg) - 1
	cpu.registers.write(reg, value)

	cpu.registers.setFlag(FLAG_Z, value == 0)
	cpu.registers.setFlag(FLAG_N, true)
	cpu.registers.setFlag(FLAG_H, (value&0x0F) == 0x0F)
}

func decR16(cpu *CPU, reg CpuRegister) {
	cpu.registers.write(reg, cpu.registers.read(reg)-1)
}

func decMR16(cpu *CPU, reg CpuRegister) {
	address := cpu.registers.read(reg)
	value := cpu.bus.readByte(address) - 1
	cpu.bus.writeByte(address, value)

	cpu.registers.setFlag(FLAG_Z, value == 0)
	cpu.registers.setFlag(FLAG_N, false)
	cpu.registers.setFlag(FLAG_H, (value&0x0F) == 0x0F)
}

func jr(cpu *CPU, cond condition) {
	e8 := readByteFromPC(cpu)

	if !checkCondition(cpu, cond) {
		return
	}

	offset := uint16(int8(e8))
	nextAddress := cpu.registers.read(R_PC) + offset
	cpu.registers.write(R_PC, nextAddress)
}

func ldhR8ToA8(cpu *CPU, src CpuRegister) {
	a8 := readByteFromPC(cpu)
	address := 0xFF00 + uint16(a8)
	cpu.bus.writeByte(address, byte(cpu.registers.read(src)))
}

func ldhA8ToR8(cpu *CPU, dest CpuRegister) {
	a8 := readByteFromPC(cpu)
	address := 0xFF00 + uint16(a8)
	cpu.registers.write(dest, uint16(cpu.bus.readByte(address)))
}

func cpN8(cpu *CPU) {
	minuend := byte(cpu.registers.read(R_A))
	subtrahend := readByteFromPC(cpu)
	cp(cpu, minuend, subtrahend)
}

func cpR8(cpu *CPU, src CpuRegister) {
	minuend := byte(cpu.registers.read(R_A))
	subtrahend := byte(cpu.registers.read(src))
	cp(cpu, minuend, subtrahend)
}

func cpMR8(cpu *CPU, src CpuRegister) {
	minuend := cpu.bus.readByte(cpu.registers.read(R_A))
	subtrahend := byte(cpu.registers.read(src))
	cp(cpu, minuend, subtrahend)
}

func cp(cpu *CPU, minuend byte, subtrahend byte) {
	cpu.registers.setFlags(
		minuend == subtrahend,
		true,
		minuend&0x0F < subtrahend*0x0F,
		minuend < subtrahend,
	)
}

func addR8(cpu *CPU, src CpuRegister) {
	a := cpu.registers.read(R_A)
	addend := cpu.registers.read(src)
	result := a + addend

	cpu.registers.write(R_A, uint16(result))

	cpu.registers.setFlags(
		result == 0,
		false,
		((a&0x0F)+(addend&0x0F)) > 0x0F,
		result > 0xFF,
	)
}

func addR16(cpu *CPU, src CpuRegister) {
	hl := cpu.registers.read(R_HL)
	addend := cpu.registers.read(src)
	result := hl + addend

	cpu.registers.write(R_HL, result)

	cpu.registers.setFlags(
		result == 0,
		false,
		((hl&0x0FFF)+(addend&0x0FFF)) > 0x0FFF,
		uint32(hl)+uint32(addend) > 0xFFFF,
	)
}

func addMR16(cpu *CPU, src CpuRegister) {
	a := cpu.registers.read(R_A)
	addend := uint16(cpu.bus.readByte(cpu.registers.read(src)))
	result := a + addend

	cpu.registers.write(R_A, uint16(result))

	cpu.registers.setFlags(
		result == 0,
		false,
		((a&0x0F)+(addend&0x0F)) > 0x0F,
		result > 0xFF,
	)
}

func addN8A(cpu *CPU) {
	a := cpu.registers.read(R_A)
	addend := uint16(readByteFromPC(cpu))
	result := a + addend

	cpu.registers.write(R_A, result)

	cpu.registers.setFlags(
		result == 0,
		false,
		((a&0x0F)+(addend&0x0F)) > 0x0F,
		result > 0xFF,
	)
}

func addN8SP(cpu *CPU) {
	sp := cpu.registers.read(R_SP)
	addend := uint16(readByteFromPC(cpu))
	result := sp + addend

	cpu.registers.write(R_HL, result)

	cpu.registers.setFlags(
		false,
		false,
		((sp&0x0FFF)+(addend&0x0FFF)) > 0x0FFF,
		uint32(sp)+uint32(addend) > 0xFFFF,
	)
}

func adcR8(cpu *CPU, src CpuRegister) {
	adc(cpu, cpu.registers.read(src))
}

func adcMR16(cpu *CPU, src CpuRegister) {
	adc(cpu, uint16(cpu.bus.readByte(cpu.registers.read(src))))
}

func adcN8(cpu *CPU) {
	adc(cpu, uint16(readByteFromPC(cpu)))
}

func adc(cpu *CPU, addend uint16) {
	a := cpu.registers.read(R_A)
	carry := BoolToUint16(cpu.registers.readFlag(FLAG_C))
	result := a + addend + carry

	cpu.registers.write(R_A, result)

	cpu.registers.setFlags(
		result == 0,
		false,
		a&0x0F+addend&0x0F+carry > 0x0F,
		result > 0xFF,
	)
}

func subR8(cpu *CPU, src CpuRegister) {
	sub(cpu, cpu.registers.read(src))
}

func subMR16(cpu *CPU, src CpuRegister) {
	sub(cpu, uint16(cpu.bus.readByte(cpu.registers.read(src))))
}

func subN8(cpu *CPU) {
	sub(cpu, uint16(readByteFromPC(cpu)))
}

func sub(cpu *CPU, subtrahend uint16) {
	a := cpu.registers.read(R_A)
	result := a - subtrahend

	cpu.registers.write(R_A, result)

	cpu.registers.setFlags(
		result == 0,
		true,
		((a & 0x0F) < (subtrahend & 0x0F)),
		a < subtrahend,
	)
}

func sbcR8(cpu *CPU, src CpuRegister) {
	sbc(cpu, cpu.registers.read(src))
}

func sbcMR16(cpu *CPU, src CpuRegister) {
	sbc(cpu, uint16(cpu.bus.readByte(cpu.registers.read(src))))
}

func sbcN8(cpu *CPU) {
	sbc(cpu, uint16(readByteFromPC(cpu)))
}

func sbc(cpu *CPU, subtrahend uint16) {
	a := cpu.registers.read(R_A)
	carry := BoolToUint16(cpu.registers.readFlag(FLAG_C))
	result := a - subtrahend - carry

	cpu.registers.write(R_A, result)

	cpu.registers.setFlags(
		result == 0,
		true,
		((a^subtrahend^uint16((byte(result)&0xFF)))&(1<<4)) != 0,
		subtrahend+carry < a,
	)
}

func andR8(cpu *CPU, src CpuRegister) {
	and(cpu, byte(cpu.registers.read(src)))
}

func andMR16(cpu *CPU, src CpuRegister) {
	and(cpu, cpu.bus.readByte(cpu.registers.read(src)))
}

func andN8(cpu *CPU) {
	and(cpu, readByteFromPC(cpu))
}

func and(cpu *CPU, comparator byte) {
	a := byte(cpu.registers.read(R_A))
	result := comparator & a

	cpu.registers.write(R_A, uint16(result))
	cpu.registers.setFlags(result == 0, false, true, false)
}

func orR8(cpu *CPU, src CpuRegister) {
	or(cpu, byte(cpu.registers.read(src)))
}

func orMR16(cpu *CPU, src CpuRegister) {
	or(cpu, cpu.bus.readByte(cpu.registers.read(src)))
}

func orN8(cpu *CPU) {
	or(cpu, readByteFromPC(cpu))
}

func or(cpu *CPU, comparator byte) {
	a := byte(cpu.registers.read(R_A))
	result := comparator | a

	cpu.registers.write(R_A, uint16(result))
	cpu.registers.setFlags(result == 0, false, false, false)
}

func rlca(cpu *CPU) {
	a := byte(cpu.registers.read(R_A))
	// get c as the most significant bit of a
	c := (a >> 7) & 1
	// rotate a 1 to the left, and put c in the rightmost position
	a = (a << 1) | c

	cpu.registers.write(R_A, uint16(a))
	cpu.registers.setFlags(false, false, false, c == 1)
}

func rrca(cpu *CPU) {
	a := byte(cpu.registers.read(R_A))

	// get c as the least significant bit of a
	c := (a & 1)
	// rotate a 1 to the right, and put c in the leftmost position
	a = (a >> 1) | c

	cpu.registers.write(R_A, uint16(a))
	cpu.registers.setFlags(false, false, false, c == 1)
}

func rla(cpu *CPU) {
	a := byte(cpu.registers.read(R_A))
	c := BoolToByte(cpu.registers.readFlag(FLAG_C))

	msb := (a >> 7) & 1
	a = (a << 1) | c

	cpu.registers.write(R_A, uint16(a))
	cpu.registers.setFlags(false, false, false, msb == 1)
}

func rra(cpu *CPU) {
	a := byte(cpu.registers.read(R_A))
	c := BoolToByte(cpu.registers.readFlag(FLAG_C))

	lsb := a & 1
	a = (a >> 1) | (c << 7)

	cpu.registers.write(R_A, uint16(a))
	cpu.registers.setFlags(false, false, false, lsb == 1)
}

// The DAA (Decimal Adjust Accumulator) instruction is used to adjust the
// accumulator register after performing a binary-coded decimal (BCD) addition
// or subtraction operation. BCD represents each decimal digit with its binary
// equivalent. The DAA instruction adjusts the result of such operations to
// ensure the accumulator contains a valid BCD number.
// i.e. this is black magic I don't fully understand, but we'll see what blargg
// has to say about it
func daa(cpu *CPU) {
	a := cpu.registers.read(R_A)
	correction := uint16(0)

	h := cpu.registers.readFlag(FLAG_H)
	if h || (a&0x0F) > 0x09 {
		correction += 0x06
	}

	c := cpu.registers.readFlag(FLAG_C)
	if c || a > 0x99 {
		correction += 0x60
		c = true
	} else {
		c = false
	}

	n := cpu.registers.readFlag(FLAG_N)
	if n {
		a -= correction
	} else {
		a += correction
	}

	cpu.registers.write(R_A, a)
	cpu.registers.setFlags(a == 0, n, false, c)
}

func cpl(cpu *CPU) {
	a := cpu.registers.read(R_A)
	r := uint16(byte(^a))
	cpu.registers.write(R_A, r)
	cpu.registers.setFlag(FLAG_N, true)
	cpu.registers.setFlag(FLAG_H, true)
}

func scf(cpu *CPU) {
	cpu.registers.setFlag(FLAG_N, false)
	cpu.registers.setFlag(FLAG_H, false)
	cpu.registers.setFlag(FLAG_C, true)
}

func ccf(cpu *CPU) {
	cpu.registers.setFlag(FLAG_N, false)
	cpu.registers.setFlag(FLAG_H, false)
	cpu.registers.setFlag(FLAG_C, !cpu.registers.readFlag(FLAG_C))
}

func push(cpu *CPU, src CpuRegister) {
	r := cpu.registers.read(src)
	hi, lo := Uint16ToBytes(r)

	decR16(cpu, R_SP)
	cpu.bus.writeByte(cpu.registers.read(R_SP), hi)
	decR16(cpu, R_SP)
	cpu.bus.writeByte(cpu.registers.read(R_SP), lo)
}

func pop(cpu *CPU, hiDest CpuRegister, loDest CpuRegister) {
	ldMR16ToR8(cpu, R_SP, loDest)
	incR16(cpu, R_SP)

	ldMR16ToR8(cpu, R_SP, hiDest)
	incR16(cpu, R_SP)
}

func jpA16(cpu *CPU, cond condition) {
	nextAddress := readWordFromPC(cpu)

	if !checkCondition(cpu, cond) {
		return
	}

	cpu.registers.write(R_PC, nextAddress)
}

func jpR16(cpu *CPU, src CpuRegister) {
	ldR16ToR16(cpu, src, R_PC)
}

func call(cpu *CPU, cond condition) {
	address := readWordFromPC(cpu)

	if !checkCondition(cpu, cond) {
		return
	}

	pc := cpu.registers.read(R_PC)
	hi, lo := Uint16ToBytes(pc)

	decR16(cpu, R_SP)
	cpu.bus.writeByte(cpu.registers.read(R_SP), hi)
	decR16(cpu, R_SP)
	cpu.bus.writeByte(cpu.registers.read(R_SP), lo)

	cpu.registers.write(R_PC, address)
}

func ret(cpu *CPU, cond condition) {
	if !checkCondition(cpu, cond) {
		return
	}

	lo := cpu.bus.readByte(cpu.registers.read(R_SP))
	incR16(cpu, R_SP)

	hi := cpu.bus.readByte(cpu.registers.read(R_SP))
	incR16(cpu, R_SP)

	cpu.registers.write(R_PC, BytesToUint16(hi, lo))
}

func rst(cpu *CPU, address uint16) {
	push(cpu, R_PC)
	cpu.registers.write(R_PC, address)
}
