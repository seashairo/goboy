package goboy

import "fmt"

type CPU struct {
}

func NewCPU() CPU {
	return CPU{}
}

func (cpu *CPU) Tick() {
	fmt.Println("CPU tick")
}
