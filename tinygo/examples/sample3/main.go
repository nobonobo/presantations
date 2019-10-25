package main

import (
	"image/color"
	"machine"
	"time"
)

const NUM_PIXELS = 10

func main() {
	machine.InitADC()
	sensor := machine.ADC{machine.A8}
	sensor.Configure()
	neo := NewNeopixelDriver(machine.GPIO{machine.NEOPIXELS}, NUM_PIXELS)
	for {
		for i := 0; i < NUM_PIXELS; i++ {
			pow := uint8(sensor.Get() / 256)
			neo.SetPixel(i, color.RGBA{R: pow, G: pow, B: pow})
		}
		neo.Show()
		time.Sleep(100 * time.Millisecond)
	}
}
