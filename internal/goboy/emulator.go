package goboy

func Emulate() {
	gameboy := NewGameBoy()
	go gameboy.Run()
	defer gameboy.Stop()

	ui := NewUI(&gameboy)
	defer ui.Destroy()

	for ui.running {
		ui.Update()
	}
}
