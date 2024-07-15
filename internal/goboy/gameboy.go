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
