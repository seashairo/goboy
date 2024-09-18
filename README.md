# Running the program

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
