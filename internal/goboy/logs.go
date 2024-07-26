package goboy

import (
	"os"
	"sync"
)

type FileSingleton struct {
	file *os.File
}

var instance *FileSingleton
var once sync.Once

func GetInstance() *FileSingleton {
	once.Do(func() {
		file, err := os.OpenFile("doctor.out", os.O_CREATE|os.O_TRUNC, 0600)

		if err != nil {
			panic(err)
		}

		instance = &FileSingleton{file: file}
	})

	return instance
}

func (fs *FileSingleton) WriteString(s string) error {
	_, err := fs.file.WriteString(s)
	return err
}

func (fs *FileSingleton) Close() error {
	return fs.file.Close()
}
