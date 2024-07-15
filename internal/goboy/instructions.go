package goboy

import "fmt"

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
	0x08: func(cpu *CPU) {
		ldR16ToA16(cpu, R_SP)
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
	0x18: func(cpu *CPU) {
		jr(cpu, true)
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
	0x20: func(cpu *CPU) {
		jr(cpu, !cpu.registers.readFlag(FLAG_Z))
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
	0x28: func(cpu *CPU) {
		jr(cpu, cpu.registers.readFlag(FLAG_Z))
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
	0x30: func(cpu *CPU) {
		jr(cpu, !cpu.registers.readFlag(FLAG_C))
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
	0x38: func(cpu *CPU) {
		jr(cpu, !cpu.registers.readFlag(FLAG_C))
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
	0xC3: func(cpu *CPU) {
		nextAddress := readWordFromPC(cpu)
		cpu.registers.write(R_PC, nextAddress)
	},
	0x31: func(cpu *CPU) {
		ldN16ToR16(cpu, R_SP)
	},
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
	0xE0: func(cpu *CPU) {
		ldhR8ToA8(cpu, R_A)
	},
	0xE2: func(cpu *CPU) {
		ldR8ToMR8(cpu, R_A, R_C)
	},
	0xEA: func(cpu *CPU) {
		ldR8ToA16(cpu, R_A)
	},
	0xEE: func(cpu *CPU) {
		xorN8(cpu)
	},
	0xF0: func(cpu *CPU) {
		ldhA8ToR8(cpu, R_A)
	},
	0xF2: func(cpu *CPU) {
		ldMR8ToR8(cpu, R_C, R_A)
	},
	0xF3: func(cpu *CPU) {
		cpu.masterInterruptEnabled = false
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
	0xFE: func(cpu *CPU) {
		cpN8(cpu)
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

func ldR8ToA16(cpu *CPU, src CpuRegister) {
	a16 := readWordFromPC(cpu)
	dest := cpu.bus.readWord(a16)
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

func jr(cpu *CPU, conditionMet bool) {
	e8 := readByteFromPC(cpu)

	if !conditionMet {
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

func cpR8(cpu *CPU, r8 CpuRegister) {
	minuend := byte(cpu.registers.read(R_A))
	subtrahend := byte(cpu.registers.read(r8))
	cp(cpu, minuend, subtrahend)
}

func cpMR8(cpu *CPU, r8 CpuRegister) {
	minuend := cpu.bus.readByte(cpu.registers.read(R_A))
	subtrahend := byte(cpu.registers.read(r8))
	cp(cpu, minuend, subtrahend)
}

func cp(cpu *CPU, minuend byte, subtrahend byte) {
	cpu.registers.setFlag(FLAG_Z, minuend == subtrahend)
	cpu.registers.setFlag(FLAG_N, true)
	cpu.registers.setFlag(FLAG_H, minuend&0x0F < subtrahend*0x0F)
	cpu.registers.setFlag(FLAG_C, minuend < subtrahend)
}
