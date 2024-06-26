package main

import (
	"github.com/seashairo/goboy/internal/goboy"
)

func main() {
	emulator := goboy.NewEmulator()
	emulator.Run()
}
