# GoBoy

Another GameBoy emulator. Probably another one called GoBoy, since I assume
other people have written GameBoy emulators in Go and used the exact same name.

Currently, you need to provide your own ROMs and modify the code to load them. I
tested with Tetris, Dr Mario, and Alleyway. I also got through the first battle
of Pokemon Red. No other games are guaranteed to work.

This emulator features

- working graphics
- working sound (although sometimes a bit poppy)
- MBC3 ROM and RAM banking

That's about it really. With this commit, it's unlikely I'll ever modify it
again. This was a fun project to learn a little bit of Go, and a little bit of
GameBoy emulation. It took longer than I expected, but I had a great time doing
it. If I was going to continue, I'd want to add:

- Loading ROMs without modifying code
- Save files (battery) and save states (state dumps)
- Real debug functionality rather than just this tile debug
  - sound waves
  - RAM viewer
  - breakpoints
- Clean the code up into properly separated and organised modules
- Better library style access to the internals to enable things like alternate
  renderers
- Attempt at WASM integration for browser play
- Color (or maybe not, it's just a basic GameBoy emulator after all)

But I'm probably not going to do that. If you want to copy my code, absolutely
do it (I certainly did), but if you want a good/feature complete emulator I'd
recommend looking elsewhere.

## Running the program

`go run cmd/goboy.go`

## But why doesn't it work?

tbh I don't know, but it's probably SDL (SDL2.dll in root required)
@see https://github.com/veandco/go-sdl2?tab=readme-ov-file#requirements

For the tests, SDL2.dll needs to be in internal/goboy - there's something weird
with Go tests and running with the file's dir as the current dir

## Where are those test files from?

https://github.com/adtennant/GameboyCPUTests

I'm sure there's some way to link to their repo properly, but I've just copied
the tests in here instead
