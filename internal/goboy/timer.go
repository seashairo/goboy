package goboy

import "fmt"

type Timer struct {
}

func NewTimer() Timer {
	return Timer{}
}

func (timer *Timer) Tick() {
	fmt.Println("Timer tick")
}
