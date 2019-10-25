package main

import (
	"image/color"
	"machine"
	"math"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

func calc(p float64) uint8 {
	i := int(64 * (math.Sin(2*math.Pi*p) + 1.0))
	if i > 255 {
		i = 255
	}
	if i < 0 {
		i = 0
	}
	return uint8(i)
}

func increase(p float64) float64 {
	p += 0.1
	if p > 1 {
		p -= float64(int(p))
	}
	return p
}

func main() {
	neo := machine.NEOPIXELS
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.New(neo)
	leds := make([]color.RGBA, 10)
	rp, gp, bp := 0.0, 1.0/3, 2.0/3
	for {
		for i := range leds {
			leds[i] = color.RGBA{
				R: calc(rp + 0.1*float64(i)),
				G: calc(gp + 0.1*float64(i)),
				B: calc(bp + 0.1*float64(i))}
		}
		ws.WriteColors(leds)
		time.Sleep(100 * time.Millisecond)
		rp = increase(rp)
		gp = increase(gp)
		bp = increase(bp)
	}
}
