package main

import (
	"device/arm"
	"image/color"
	"machine"

	"github.com/tinygo-org/drivers/ws2812"
)

// NeopixelDriver represents a connection to a NeoPixel
type NeopixelDriver struct {
	pin        machine.GPIO
	device     ws2812.Device
	pixelCount int
	pixels     []color.RGBA
}

// NewNeopixelDriver returns a new NeopixelDriver
func NewNeopixelDriver(pin machine.GPIO, pixelCount int) *NeopixelDriver {
	pin.Configure(machine.GPIOConfig{Mode: machine.GPIO_OUTPUT})
	neo := &NeopixelDriver{
		pin:        pin,
		device:     ws2812.New(pin),
		pixelCount: pixelCount,
		pixels:     make([]color.RGBA, pixelCount),
	}

	return neo
}

// Pin returns the Driver's pin
func (neo *NeopixelDriver) Pin() machine.GPIO { return neo.pin }

// Show activates all the Neopixels in the strip
func (neo *NeopixelDriver) Show() {
	mask := arm.DisableInterrupts()
	defer arm.EnableInterrupts(mask)
	neo.device.WriteColors(neo.pixels)
}

// SetPixel sets the color of one specific Neopixel in the strip
func (neo *NeopixelDriver) SetPixel(pix int, color color.RGBA) {
	neo.pixels[pix] = color
}
