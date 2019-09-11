package main

import (
	"image/color"
	"machine"
	"time"
)

const NUM_PIXELS = 10

func main() {
	neo := NewNeopixelDriver(machine.GPIO{machine.NEOPIXELS}, NUM_PIXELS)
	n := 0
	var pows = []uint8{255, 127, 63, 31, 15, 7, 3, 1, 0, 0}
	for {
		for i := 0; i < NUM_PIXELS; i++ {
			pow := pows[NUM_PIXELS-1-i]
			neo.SetPixel((n+i)%NUM_PIXELS, color.RGBA{R: pow, G: pow, B: pow})
		}
		neo.Show()
		time.Sleep(100 * time.Millisecond)
		n = (n + 1) % NUM_PIXELS
	}
}
