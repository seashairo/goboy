package goboy

import (
	"fmt"
	"time"
)

const ROM_PATH = "./data/roms/tetris.gb"

type GameBoy struct {
	running bool
	paused  bool
	ticks   uint64

	cpu   CPU
	ppu   PPU
	timer Timer
	bus   Bus
}

func NewGameBoy() GameBoy {
	bus := NewBus(ROM_PATH)
	cpu := NewCPU(bus)

	return GameBoy{
		running: false,
		paused:  false,
		ticks:   0,
		cpu:     cpu,
		ppu:     NewPPU(),
		timer:   NewTimer(),
		bus:     bus,
	}
}

func (gameboy *GameBoy) Run() {
	debugInstructionCount()

	gameboy.running = true

	for gameboy.running {
		if gameboy.paused {
			time.Sleep(16 * time.Millisecond)
			continue
		}

		gameboy.cpu.Tick()
		gameboy.timer.Tick()
		gameboy.ppu.Tick()

		gameboy.ticks += 1
	}

	fmt.Println("Terminating...")
}

func debugInstructionCount() {
	instCount := 0
	outStr := "Unimplemented opcodes: "
	for k, v := range instructions {
		if v == nil {
			instCount++
			outStr += fmt.Sprintf("0x%2.2X, ", k)
		}
	}

	fmt.Printf("%d instructions not implemented out of %d\n", instCount, 0x100)
	fmt.Println(outStr)
}
