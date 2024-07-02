package goboy

import (
	"fmt"
)

type CPU struct {
	registers Registers
}

func NewCPU() CPU {
	return CPU{
		registers: NewRegisters(),
	}
}

func (cpu *CPU) Tick() {
	fmt.Println("CPU tick")
}
