package goboy

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const ROM_PATH = "./data/roms/blargg/01-special.gb"

type GameBoy struct {
	running bool
	paused  bool
	cycles  uint64

	cpu   CPU
	ppu   PPU
	timer *Timer
	bus   Bus
}

func NewGameBoy() GameBoy {
	timer := NewTimer()
	bus := NewBus(ROM_PATH, &timer)
	cpu := NewCPU(&bus, &timer)
	ppu := NewPPU(&bus)

	return GameBoy{
		running: false,
		paused:  false,
		cycles:  0,

		cpu:   cpu,
		timer: &timer,
		bus:   bus,
		ppu:   ppu,
	}
}

func (gameboy *GameBoy) Run() {
	gameboy.running = true
	stepping := false
	input := bufio.NewReader(os.Stdin)
	gameboy.cpu.debugPrint()

	for gameboy.running {
		if gameboy.paused {
			time.Sleep(16 * time.Millisecond)
			continue
		}

		if stepping {
			in, _ := input.ReadBytes('\n')
			trimmed := strings.TrimSpace(string(in))

			if trimmed == "continue" {
				stepping = false
			} else if trimmed == "dump" {
				out := ""

				for i := uint16(0xC000); i < 0xC800; i++ {
					out += fmt.Sprintf("%2.2X ", gameboy.cpu.bus.readByte(i))
					if i%16 == 15 {
						out += "\n"
					}
				}
				fmt.Println(out)
			}
		}

		gameboy.cpu.Tick()
		gameboy.ppu.Tick()

		gameboy.cycles++
	}

	fmt.Println("GameBoy terminating")
}

func (gameboy *GameBoy) Stop() {
	gameboy.running = false
}
