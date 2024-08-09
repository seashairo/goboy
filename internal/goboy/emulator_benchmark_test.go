package goboy

import (
	"os"
	"path"
	"runtime"
	"testing"
	"time"
)

func BenchmarkEmulate(b *testing.B) {
	gameboy := NewGameBoy()
	go gameboy.Run()

	start := time.Now().UnixMilli()
	for {
		if time.Now().UnixMilli()-start > 10000 {
			break
		}
	}

	gameboy.Stop()
}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..") // change to suit test file location
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
