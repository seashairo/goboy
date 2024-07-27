package goboy

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var ROMS = []string{
	"./data/roms/blargg/01-special.gb",
	"./data/roms/blargg/02-interrupts.gb",
	"./data/roms/blargg/03-op sp,hl.gb",
	"./data/roms/blargg/04-op r,imm.gb",
	"./data/roms/blargg/05-op rp.gb",
	"./data/roms/blargg/06-ld r,r.gb",
	"./data/roms/blargg/07-jr,jp,call,ret,rst.gb",
	"./data/roms/blargg/08-misc instrs.gb",
	"./data/roms/blargg/09-op r,r.gb",
	"./data/roms/blargg/10-bit ops.gb",
	"./data/roms/blargg/11-op a,(hl).gb",
}

var ROM_PATH = ROMS[10]

type GameBoy struct {
	running bool
	paused  bool

	cpu   CPU
	ppu   PPU
	timer *Timer
	bus   Bus
}

func NewGameBoy() GameBoy {
	timer := NewTimer()
	bus := NewBus(ROM_PATH, &timer)
	cpu := NewCPU(bus, &timer)
	ppu := NewPPU()

	return GameBoy{
		running: false,
		paused:  false,

		cpu:   cpu,
		ppu:   ppu,
		timer: &timer,
		bus:   bus,
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
	}

	fmt.Println("Terminating...")
}
