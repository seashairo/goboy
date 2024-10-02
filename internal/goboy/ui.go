package goboy

// typedef unsigned char Uint8;
// void AudioCallback(void *userdata, Uint8 *stream, int len);
import "C"

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

var PALLETTE = [4]uint32{0xFFFFFFFF, 0xFFAAAAAA, 0xFF555555, 0xFF000000}

var scale = int32(2)

const TILE_HEIGHT = 8
const TILE_WIDTH = 8
const TILES_X = 16
const TILES_Y = 24

type UI struct {
	running bool
	gameboy *GameBoy

	lcdWindow   *sdl.Window
	lcdRenderer *sdl.Renderer
	lcdTexture  *sdl.Texture
	lcdSurface  *sdl.Surface

	previousFrame uint32

	tileDebugWindow   *sdl.Window
	tileDebugRenderer *sdl.Renderer
	tileDebugTexture  *sdl.Texture
	tileDebugSurface  *sdl.Surface

	// audioDebugWindow   *sdl.Window
	// audioDebugRenderer *sdl.Renderer
	// audioDebugTexture  *sdl.Texture
	// audioDebugSurface  *sdl.Surface
	// lastAudioSamples   []int16

	audioDeviceId sdl.AudioDeviceID
	audioBuffer   []int16
}

func NewUI(gameboy *GameBoy) *UI {
	ui := &UI{
		running:       true,
		gameboy:       gameboy,
		previousFrame: 0,
	}

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}

	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		panic(err)
	}

	ui.initLcd()
	ui.initTileDebug()
	// ui.initAudioDebug()
	ui.initAudio()

	return ui
}

func (ui *UI) Update() {
	ui.handleEvents()
	ui.updateTileDebugWindow()
	ui.updateLcdWindow()
	// ui.updateAudioDebugWindow()
}

const (
	bufferSize = 1024 // Buffer size
)

func (ui *UI) queueAudio(left int16, right int16) {
	if len(ui.audioBuffer) <= sampleRate*2 {
		ui.audioBuffer = append(ui.audioBuffer, left, right)

		// ui.lastAudioSamples = append(ui.lastAudioSamples, left)
		// if len(ui.lastAudioSamples) > sampleRate {
		// 	ui.lastAudioSamples = ui.lastAudioSamples[1:]
		// }
	}

	// 0.25s worth of audio queued up
	if sdl.GetQueuedAudioSize(ui.audioDeviceId) > sampleRate/4 {
		return
	}

	if len(ui.audioBuffer) < bufferSize*2 {
		return
	}

	byteBuffer := unsafe.Slice((*byte)(unsafe.Pointer(&ui.audioBuffer[0])), len(ui.audioBuffer)*2)

	if err := sdl.QueueAudio(ui.audioDeviceId, byteBuffer); err != nil {
		panic(err)
	}

	ui.audioBuffer = ui.audioBuffer[:0]
}

func (ui *UI) initAudio() {
	spec := sdl.AudioSpec{
		Freq:     sampleRate,
		Format:   sdl.AUDIO_S16SYS, // Signed 16-bit samples in system byte order
		Channels: 2,                // Stereo
		Samples:  bufferSize,       // Buffer size (affects the latency)
	}

	audioDeviceId, err := sdl.OpenAudioDevice("", false, &spec, nil, 0)
	if err != nil {
		panic(err)
	}
	ui.audioDeviceId = audioDeviceId

	ui.audioBuffer = make([]int16, bufferSize*2)

	sdl.PauseAudioDevice(audioDeviceId, false)

	ui.gameboy.RegisterAudioCallback(func(left int16, right int16) {
		ui.queueAudio(left, right)
	})
}

// @see https://gbdev.io/pandocs/Tile_Data.html
func (ui *UI) displayTile(tileNum uint16, xDraw int32, yDraw int32) {
	// Each tile occupies 16 bytes
	for y := int32(0); y < 16; y += 2 {
		// Where each line is represented by 2 bytes
		b1 := ui.gameboy.bus.readByte(VIDEO_RAM_START + tileNum*TILES_X + uint16(y))
		b2 := ui.gameboy.bus.readByte(VIDEO_RAM_START + tileNum*TILES_X + uint16(y) + 1)

		for bit := 7; bit >= 0; bit-- {
			// For each line, the first byte specifies the least significant bit of
			// the color ID of each pixel, and the second byte specifies the most
			// significant bit. In both bytes, bit 7 represents the leftmost pixel,
			// and bit 0 the rightmost.
			lo := (b1 & (1 << bit)) >> bit
			hi := ((b2 & (1 << bit)) >> bit) << 1
			color := hi | lo

			rect := sdl.Rect{
				X: xDraw + (7-int32(bit))*scale,
				Y: yDraw + (y/2)*scale,
				W: scale,
				H: scale,
			}

			ui.tileDebugSurface.FillRect(&rect, PALLETTE[color])
		}
	}
}

func (ui *UI) updateTileDebugWindow() {
	surface := ui.tileDebugSurface

	rect := sdl.Rect{X: 0, Y: 0, W: surface.W, H: surface.H}
	surface.FillRect(&rect, 0xFFFF0000)

	xDraw, yDraw := int32(0), int32(0)
	tileNum := uint16(0)

	for y := int32(0); y < TILES_Y; y++ {
		for x := int32(0); x < TILES_X; x++ {
			ui.displayTile(tileNum, xDraw+(x*scale), yDraw+(y*scale))
			xDraw += TILE_WIDTH * scale
			tileNum++
		}
		yDraw += TILE_HEIGHT * scale
		xDraw = 0
	}

	pixels := surface.Pixels()
	ui.tileDebugTexture.Update(&rect, unsafe.Pointer(&(pixels[0])), int(surface.Pitch))
	ui.tileDebugRenderer.Copy(ui.tileDebugTexture, nil, nil)
	ui.tileDebugRenderer.Present()
}

func (ui *UI) initTileDebug() {
	tileDebugWidth := (TILES_X * TILE_WIDTH * scale) + (TILES_X * scale) - scale
	tileDebugHeight := (TILES_Y * TILE_HEIGHT * scale) + (TILES_Y * scale) - scale
	tileDebugWindow, tileDebugRenderer, err := sdl.CreateWindowAndRenderer(tileDebugWidth, tileDebugHeight, 0)
	if err != nil {
		panic(err)
	}
	tileDebugWindow.SetTitle("Tile Debug")

	tileDebugTexture, err := tileDebugRenderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, tileDebugWidth, tileDebugHeight)
	if err != nil {
		panic(err)
	}

	x, _ := ui.lcdWindow.GetPosition()
	_, h := ui.lcdWindow.GetSize()

	tileDebugWindow.SetPosition(x, h+64)
	tileDebugSurface, err := tileDebugWindow.GetSurface()
	if err != nil {
		panic(err)
	}
	tileDebugSurface.FillRect(nil, 0)

	ui.tileDebugWindow = tileDebugWindow
	ui.tileDebugRenderer = tileDebugRenderer
	ui.tileDebugSurface = tileDebugSurface
	ui.tileDebugTexture = tileDebugTexture
}

// func (ui *UI) updateAudioDebugWindow() {
// 	surface := ui.audioDebugSurface

// 	surfaceRect := sdl.Rect{X: 0, Y: 0, W: surface.W, H: surface.H}
// 	surface.FillRect(&surfaceRect, 0xFF000000)

// 	pixels := surface.Pixels()
// 	ui.audioDebugTexture.Update(&surfaceRect, unsafe.Pointer(&(pixels[0])), int(surface.Pitch))
// 	ui.audioDebugRenderer.Copy(ui.audioDebugTexture, nil, nil)

// 	ui.audioDebugRenderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)

// 	for index, element := range ui.lastAudioSamples {
// 		if index%8 != 0 {
// 			continue
// 		}
// 		ui.audioDebugRenderer.DrawPoint(int32(index%8), int32(100-element/50))
// 	}

// 	ui.audioDebugRenderer.Present()
// }

// func (ui *UI) initAudioDebug() {
// 	ui.lastAudioSamples = make([]int16, sampleRate)

// 	audioDebugWindow, audioDebugRenderer, err := sdl.CreateWindowAndRenderer(1000, 100, 0)
// 	if err != nil {
// 		panic(err)
// 	}
// 	audioDebugWindow.SetTitle("Audio Debug")

// 	audioDebugTexture, err := audioDebugRenderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, 1000, 100)
// 	if err != nil {
// 		panic(err)
// 	}

// 	x, y := ui.lcdWindow.GetPosition()

// 	audioDebugWindow.SetPosition(x+LCD_WIDTH*scale, y)
// 	audioDebugSurface, err := audioDebugWindow.GetSurface()
// 	if err != nil {
// 		panic(err)
// 	}
// 	audioDebugSurface.FillRect(nil, 0)

// 	ui.audioDebugWindow = audioDebugWindow
// 	ui.audioDebugRenderer = audioDebugRenderer
// 	ui.audioDebugSurface = audioDebugSurface
// 	ui.audioDebugTexture = audioDebugTexture
// }

func (ui *UI) updateLcdWindow() {
	if ui.previousFrame == ui.gameboy.ppu.currentFrame {
		return
	}

	surface := ui.lcdSurface

	surfaceRect := sdl.Rect{X: 0, Y: 0, W: surface.W, H: surface.H}
	surface.FillRect(&surfaceRect, 0xFFFF0000)

	for lineNum := int32(0); lineNum < LCD_HEIGHT; lineNum++ {
		for x := int32(0); x < LCD_WIDTH; x++ {
			rect := sdl.Rect{
				X: x * scale,
				Y: lineNum * scale,
				W: scale,
				H: scale,
			}

			index := x + (lineNum * LCD_WIDTH)
			pixel := ui.gameboy.ppu.videoBuffer[index]

			surface.FillRect(&rect, pixel)
		}
	}

	pixels := surface.Pixels()
	ui.lcdTexture.Update(&surfaceRect, unsafe.Pointer(&(pixels[0])), int(surface.Pitch))
	ui.lcdRenderer.Copy(ui.lcdTexture, nil, nil)

	// ui.lcdRenderer.SetDrawColor(0, 0, 255, 255)
	// windowRect2 := sdl.Rect{
	// 	X: int32(ui.gameboy.bus.readByte(LCD_WX)) * scale,
	// 	Y: int32(ui.gameboy.bus.readByte(LCD_WY)) * scale,
	// 	W: 256 * scale,
	// 	H: 256 * scale,
	// }
	// ui.lcdRenderer.DrawRect(&windowRect2)

	ui.lcdRenderer.Present()
}

func (ui *UI) initLcd() {
	lcdWidth := LCD_WIDTH * scale
	lcdHeight := LCD_HEIGHT * scale

	lcdWindow, lcdRenderer, err := sdl.CreateWindowAndRenderer(lcdWidth, lcdHeight, 0)
	if err != nil {
		panic(err)
	}
	lcdWindow.SetTitle("GoBoy")
	lcdWindow.SetPosition(0, 32)

	lcdSurface, err := lcdWindow.GetSurface()
	if err != nil {
		panic(err)
	}
	lcdSurface.FillRect(nil, 0)

	lcdTexture, err := lcdRenderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, lcdWidth, lcdHeight)
	if err != nil {
		panic(err)
	}

	ui.lcdWindow = lcdWindow
	ui.lcdRenderer = lcdRenderer
	ui.lcdSurface = lcdSurface
	ui.lcdTexture = lcdTexture
}

func (ui *UI) handleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case sdl.KeyboardEvent:
			var b Button = 255
			switch t.Keysym.Sym {
			case sdl.K_z:
				b = JOYPAD_A
			case sdl.K_x:
				b = JOYPAD_B
			case sdl.K_BACKSPACE:
				b = JOYPAD_SELECT
			case sdl.K_RETURN:
				b = JOYPAD_START
			case sdl.K_RIGHT:
				b = JOYPAD_RIGHT
			case sdl.K_LEFT:
				b = JOYPAD_LEFT
			case sdl.K_UP:
				b = JOYPAD_UP
			case sdl.K_DOWN:
				b = JOYPAD_DOWN
			}

			if b == 255 {
				continue
			}

			if t.Type == sdl.KEYDOWN {
				ui.gameboy.Press(b)
			} else if t.Type == sdl.KEYUP {
				ui.gameboy.Release(b)
			}
		case sdl.WindowEvent:
			if t.Event == sdl.WINDOWEVENT_CLOSE {
				ui.running = false
			}
		case sdl.QuitEvent:
			println("Quit")
			ui.running = false
		}
	}
}

func (ui *UI) Destroy() {
	ui.lcdWindow.Destroy()
	ui.tileDebugWindow.Destroy()

	sdl.PauseAudioDevice(ui.audioDeviceId, false)
	sdl.CloseAudioDevice(ui.audioDeviceId)

	sdl.Quit()
}
