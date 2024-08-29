package goboy

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// const ROM_PATH = "./data/roms/drmario.gb"

const ROM_PATH = "./data/roms/dmg-acid2.gb"

// const ROM_PATH = "./data/roms/blargg/01-special.gb"

// const ROM_PATH = "./data/roms/blargg/02-interrupts.gb"

// const ROM_PATH = "./data/roms/blargg/03-op sp,hl.gb"

// const ROM_PATH = "./data/roms/blargg/04-op r,imm.gb"

// const ROM_PATH = "./data/roms/blargg/05-op rp.gb"

// const ROM_PATH = "./data/roms/blargg/06-ld r,r.gb"

// const ROM_PATH = "./data/roms/blargg/07-jr,jp,call,ret,rst.gb"

// const ROM_PATH = "./data/roms/blargg/08-misc instrs.gb"

// const ROM_PATH = "./data/roms/blargg/09-op r,r.gb"

// const ROM_PATH = "./data/roms/blargg/10-bit ops.gb"

// const ROM_PATH = "./data/roms/blargg/11-op a,(hl).gb"

type GameBoy struct {
	running bool
	paused  bool
	cycles  uint64

	cpu   *CPU
	ppu   *PPU
	timer *Timer
	bus   *Bus
}

func NewGameBoy() *GameBoy {
	bus := &Bus{}

	// Initialize all the Game Boy hardware
	timer := NewTimer()
	cpu := NewCPU(bus, timer)
	ppu := NewPPU(bus)

	cartridge := LoadCartridge(ROM_PATH)
	wram := NewRAM(8192, WORK_RAM_START)
	hram := NewRAM(127, HIGH_RAM_START)
	interruptFlagsRegister := NewInterruptRegister(0)
	io := NewIO(bus, timer, interruptFlagsRegister)
	interruptEnableRegister := NewInterruptRegister(0)
	// And then put it on the bus so everything knows what it has access to
	bus.Init(cartridge, ppu, wram, hram, io, interruptEnableRegister)

	return &GameBoy{
		running: false,
		paused:  false,
		cycles:  0,

		cpu:   cpu,
		timer: timer,
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

		gameboy.cycles++
	}

	fmt.Println("GameBoy terminating")
}

func (gameboy *GameBoy) Stop() {
	gameboy.running = false
}
