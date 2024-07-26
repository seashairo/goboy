package goboy

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const ROM_PATH = "./data/roms/blargg/07-jr,jp,call,ret,rst.gb"

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
	ppu := NewPPU()
	timer := NewTimer()

	return GameBoy{
		running: false,
		paused:  false,
		ticks:   0,
		cpu:     cpu,
		ppu:     ppu,
		timer:   timer,
		bus:     bus,
	}
}

func (gameboy *GameBoy) Run() {
	gameboy.running = true
	stepping := false
	input := bufio.NewReader(os.Stdin)

	for gameboy.running {
		if gameboy.paused {
			time.Sleep(16 * time.Millisecond)
			continue
		}

		// if gameboy.cpu.registers.read(R_PC) == 0xC659 {
		if gameboy.ticks == 1070852 {
			fmt.Println("\nEntering step mode")
			// stepping = true
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
		gameboy.timer.Tick()
		gameboy.ppu.Tick()

		gameboy.ticks += 1

		if gameboy.ticks > 1000000 {
			break
		}
	}

	fmt.Println("Terminating...")
}
