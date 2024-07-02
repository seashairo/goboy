package goboy

import (
	"fmt"
	"time"
)

const ROM_PATH = "./data/roms/pokemon_red.gb"

type Emulator struct {
	running bool
	paused  bool
	ticks   uint64

	cpu   CPU
	ppu   PPU
	timer Timer
	bus   Bus
}

func NewEmulator() Emulator {
	return Emulator{
		running: false,
		paused:  false,
		ticks:   0,
		cpu:     NewCPU(),
		ppu:     NewPPU(),
		timer:   NewTimer(),
		bus:     NewBus(ROM_PATH),
	}
}

func (emulator *Emulator) Run() {
	emulator.running = true

	for emulator.running {
		if emulator.paused {
			time.Sleep(16 * time.Millisecond)
			continue
		}

		emulator.cpu.Tick()
		emulator.timer.Tick()
		emulator.ppu.Tick()
		emulator.ticks += 1

		if emulator.ticks == 3 {
			emulator.running = false
		}
	}

	fmt.Println("Terminating...")
}
