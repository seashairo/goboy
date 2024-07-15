package main

import (
	"github.com/seashairo/goboy/internal/goboy"
)

func main() {
	gameboy := goboy.NewGameBoy()
	gameboy.Run()
}
