package goboy

import (
	"fmt"
	"time"
)

type FPSTimer struct {
	name           string
	frameCount     int
	startTime      time.Time
	rateLimit      time.Duration
	frameStartTime time.Time
}

func NewFPSTimer(name string, fpsLimit int) *FPSTimer {
	var frameDuration time.Duration
	if fpsLimit > 0 {
		frameDuration = time.Second / time.Duration(fpsLimit)
	}
	return &FPSTimer{
		name:      name,
		startTime: time.Now(),
		rateLimit: frameDuration,
	}
}

func (t *FPSTimer) FrameStart() {
	if t.rateLimit > 0 {
		if !t.frameStartTime.IsZero() {
			elapsed := time.Since(t.frameStartTime)
			if elapsed < t.rateLimit {
				time.Sleep(t.rateLimit - elapsed)
			}
		}
	}

	t.frameStartTime = time.Now()
}

func (t *FPSTimer) FrameEnd() {
	t.frameCount++

	if time.Since(t.startTime) >= time.Second {
		fmt.Printf("%s: %d\n", t.name, t.frameCount)

		t.startTime = time.Now()
		t.frameCount = 0
	}
}
