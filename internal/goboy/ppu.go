package goboy

import "fmt"

type PPU struct {
}

func NewPPU() PPU {
	return PPU{}
}

func (ppu *PPU) Tick() {
	fmt.Println("PPU tick")
}
