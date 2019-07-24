package main

import (
	"image/color"
	"machine"
	"time"
)

const NUM_PIXELS = 10

func observer(event chan int) {
	btnA := machine.GPIO{machine.BUTTONA}
	btnA.Configure(machine.GPIOConfig{Mode: machine.GPIO_INPUT_PULLDOWN})
	btnB := machine.GPIO{machine.BUTTONB}
	btnB.Configure(machine.GPIOConfig{Mode: machine.GPIO_INPUT_PULLDOWN})
	next := make([]bool, 2)
	prev := make([]bool, 2)
	prev[0], prev[1] = btnA.Get(), btnB.Get()
	for {
		time.Sleep(20 * time.Millisecond)
		next[0], next[1] = btnA.Get(), btnB.Get()
		for i := range next {
			if !prev[i] && next[i] {
				event <- i
			}
		}
		prev[0], prev[1] = next[0], next[1]
	}
}

func renderer(ch chan *NeopixelDriver) {
	for {
		v := <-ch
		v.Show()
	}
}

func main() {
	neo := NewNeopixelDriver(machine.GPIO{machine.NEOPIXELS}, NUM_PIXELS)
	ch := make(chan *NeopixelDriver)
	go renderer(ch)
	events := make(chan int)
	go observer(events)
	prev, next := 0, 0
	neo.SetPixel(next, color.RGBA{R: 10, G: 10, B: 10})
	ch <- neo
	for {
		select {
		case v := <-events:
			switch v {
			case 0:
				println("buttonA pressed")
				next++
				if next >= NUM_PIXELS {
					next = 0
				}
			case 1:
				println("buttonA pressed")
				next--
				if next < 0 {
					next = NUM_PIXELS - 1
				}
			default:
				continue
			}
		}
		neo.SetPixel(prev, color.RGBA{R: 0, G: 0, B: 0})
		neo.SetPixel(next, color.RGBA{R: 10, G: 10, B: 10})
		ch <- neo
		prev = next
	}
}
