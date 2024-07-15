package goboy

import (
	"fmt"
)

type CPU struct {
	registers CpuRegisters
	bus       Bus

	halted                 bool
	masterInterruptEnabled bool
}

func NewCPU(bus Bus) CPU {
	return CPU{
		registers:              NewCpuRegisters(),
		bus:                    bus,
		halted:                 false,
		masterInterruptEnabled: false,
	}
}

func (cpu *CPU) Tick() {

	if cpu.halted {
		return
	}

	fmt.Println("")
	cpu.registers.debugPrint()

	currentOpcode := cpu.fetchNextOpcode()

	fmt.Printf("Opcode: 0x%2.2X\n", currentOpcode)

	instruction := fetchInstruction(currentOpcode)
	instruction(cpu)
}

func (cpu *CPU) fetchNextOpcode() byte {
	pc := cpu.registers.read(R_PC)
	opcode := cpu.bus.readByte(pc)
	cpu.registers.write(R_PC, pc+1)

	return opcode
}
