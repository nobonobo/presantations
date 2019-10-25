package main

import (
	"machine"
	"time"
)

var LED machine.Pin

func init() {
	LED = machine.Pin(6) // <- ピン番号（ターゲットに合わせる）
	LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
}
func main() {
	for {
		LED.High()
		time.Sleep(1 * time.Second)
		LED.Low()
		time.Sleep(1 * time.Second)
	}
}
