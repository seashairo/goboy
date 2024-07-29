package goboy

import "github.com/veandco/go-sdl2/sdl"

type UI struct {
	running bool

	window  *sdl.Window
	surface *sdl.Surface
}

func NewUI() UI {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(
		"test",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		800,
		600,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		panic(err)
	}

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	return UI{
		running: true,
		window:  window,
		surface: surface,
	}
}

func (ui *UI) Update() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case sdl.QuitEvent:
			println("Quit")
			ui.running = false
			break
		}
	}
}

func (ui *UI) Destroy() {
	ui.window.Destroy()
	sdl.Quit()
}
