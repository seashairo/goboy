package goboy

import (
	"fmt"
	"os"
)

type CPU struct {
	registers CpuRegisters
	bus       Bus

	halted                  bool
	interruptMasterEnabled  bool
	enablingInterruptMaster bool
}

func NewCPU(bus Bus) CPU {
	return CPU{
		registers:               NewCpuRegisters(),
		bus:                     bus,
		halted:                  false,
		interruptMasterEnabled:  false,
		enablingInterruptMaster: false,
	}
}

func (cpu *CPU) Tick() {
	if cpu.halted {
		if cpu.bus.io.interrupts.readByte() != 0 {
			cpu.halted = false
		}
	} else {
		// fmt.Println("")
		cpu.debugPrint()

		currentOpcode := cpu.fetchNextOpcode()

		// fmt.Printf("Opcode: 0x%2.2X\n", currentOpcode)

		instruction := fetchInstruction(currentOpcode)
		instruction(cpu)
	}

	if cpu.interruptMasterEnabled {
		cpu.handleInterrupts()
		cpu.enablingInterruptMaster = false
	}

	if cpu.enablingInterruptMaster {
		cpu.interruptMasterEnabled = true
	}
}

func (cpu *CPU) fetchNextOpcode() byte {
	pc := cpu.registers.read(R_PC)
	opcode := cpu.bus.readByte(pc)
	cpu.registers.write(R_PC, pc+1)

	return opcode
}

func (cpu *CPU) debugPrint() {
	r := cpu.registers

	out := fmt.Sprintf(
		"A:%2.2X F:%2.2X B:%2.2X C:%2.2X D:%2.2X E:%2.2X H:%2.2X L:%2.2X SP:%4.4X PC:%4.4X PCMEM:%2.2X,%2.2X,%2.2X,%2.2X\n",
		r.a,
		r.f,
		r.b,
		r.c,
		r.d,
		r.e,
		r.h,
		r.l,
		r.sp,
		r.pc,
		cpu.bus.readByte(r.pc),
		cpu.bus.readByte(r.pc+1),
		cpu.bus.readByte(r.pc+2),
		cpu.bus.readByte(r.pc+3),
	)
	// fmt.Print(out)

	// fmt.Printf(
	// 	"Flags - Z: %t, N: %t, H: %t, C: %t\n",
	// 	r.readFlag(FLAG_Z),
	// 	r.readFlag(FLAG_N),
	// 	r.readFlag(FLAG_H),
	// 	r.readFlag(FLAG_C),
	// )

	f, err := os.OpenFile("doctor.out", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(string(out)); err != nil {
		panic(err)
	}
}
