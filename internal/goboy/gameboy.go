package goboy

import (
	"fmt"
	"time"
)

const DEBUG = false

// const ROM_PATH = "./data/roms/alleyway.gb"

// const ROM_PATH = "./data/roms/drmario.gb"

// const ROM_PATH = "./data/roms/dmg-acid2.gb"

// const ROM_PATH = "./data/roms/blargg/instr_timing.gb"

// const ROM_PATH = "./data/roms/blargg/01-special.gb"

const ROM_PATH = "./data/roms/blargg/02-interrupts.gb"

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
	*Joypad

	running bool
	paused  bool
	cycles  uint64

	cpu   *CPU
	ppu   *PPU
	timer *Timer
	bus   MemoryBusser
	io    *IO
}

func NewGameBoy() *GameBoy {
	gameboy := &GameBoy{
		running: false,
		paused:  false,
		cycles:  0,
	}

	bus := &Bus{}
	lcd := NewLCD(gameboy, bus)

	// Initialize all the Game Boy hardware
	cpu := NewCPU(gameboy, bus)
	ppu := NewPPU(gameboy, bus, lcd)

	cartridge := LoadCartridge(ROM_PATH)
	wram := NewRAM(8192, WORK_RAM_START)
	hram := NewRAM(127, HIGH_RAM_START)

	interruptFlagsRegister := NewInterruptRegister(0)
	joypad := NewJoypad(gameboy, bus)
	timer := NewTimer(gameboy)

	io := NewIO(gameboy, bus, timer, interruptFlagsRegister, lcd, joypad)
	interruptEnableRegister := NewInterruptRegister(0)
	// And then put it on the bus so everything knows what it has access to
	bus.Init(cartridge, ppu, wram, hram, io, interruptEnableRegister)

	gameboy.cpu = cpu
	gameboy.timer = timer
	gameboy.bus = bus
	gameboy.ppu = ppu
	gameboy.io = io
	gameboy.Joypad = joypad

	return gameboy
}

func (gameboy *GameBoy) Run() {
	gameboy.running = true
	gameboy.cpu.debugPrint()

	for gameboy.running {
		if gameboy.paused {
			time.Sleep(16 * time.Millisecond)
			continue
		}

		gameboy.cpu.Tick()
		gameboy.cycles++
	}

	fmt.Println("GameBoy terminating")
}

func (gameboy *GameBoy) RequestInterrupt(kind InterruptKind) {
	gameboy.io.interrupts.SetInterrupt(kind, true)
}

func (gameboy *GameBoy) ClearInterrupt(kind InterruptKind) {
	gameboy.io.interrupts.SetInterrupt(kind, false)
}

func (gameboy *GameBoy) Cycle(mCycles int) {
	tCycles := mCycles * 4
	for i := 0; i < tCycles; i++ {
		gameboy.timer.Tick()
		gameboy.io.dma.Tick()
		gameboy.ppu.Tick()
		gameboy.io.serial.Tick()
	}
}

func (gameboy *GameBoy) Stop() {
	gameboy.running = false
}
