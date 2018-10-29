package main

import "C"
import (
	"github.com/janne/bcm2835"
	"time"
)

const PIN = bcm2835.Pin11

const (
	BCM2835_GPIO_FSEL_INPT = 0
	BCM2835_GPIO_FSEL_OUTP = 1
)

func delay(microSec time.Duration) {
	time.Sleep(microSec * time.Microsecond)
}

func up() {
	bcm2835.GpioFsel(PIN, BCM2835_GPIO_FSEL_INPT)
	bcm2835.GpioSet(PIN)
}

func down() {
	bcm2835.GpioFsel(PIN, BCM2835_GPIO_FSEL_OUTP)
	bcm2835.GpioClr(PIN)
}

func send(b byte) {
	for bit := 0; bit < 8; bit++ {
		down()

		if b&1 != 0 {
			delay(8)
		} else {
			delay(44)
		}

		up()

		if b&1 != 0 {
			delay(44)
		} else {
			delay(8)
		}

		b = b >> 1
	}

	delay(1)
}

func write(b byte) {
	for bit := 0; bit < 8; bit++ {

		down()

		if b&1 != 0 {
			delay(50)
		} else {
			delay(1)
		}

		up()
		delay(10000)

		b = b >> 1
	}
}

func sendR() {
	down()
	delay(1024)
	up()
	delay(1)
}

func waitP() int {
	count := 10000

	for bcm2835.GpioLev(PIN) != 0 {
		count--

		if count < 0 {
			return 0
		}
	}
	for bcm2835.GpioLev(PIN) == 0 {
		delay(1)
	}

	delay(1)

	return 1
}

func recv() byte {
	var b byte = 0

	for bit := 0; bit < 8; bit++ {

		down()
		delay(8)

		up()
		delay(8)

		b >>= 1

		if bcm2835.GpioLev(PIN) != 0 {
			b |= 0x80
		}

		delay(32)
	}

	return b
}

func main() {
	err := bcm2835.Init()

	if err != nil {
		panic(err)
	}

	bcm2835.GpioFsel(PIN, bcm2835.Input)
	bcm2835.GpioSetPud(PIN, 0)
	bcm2835.GpioSet(PIN)

	defer bcm2835.Close()

	for {

		delay(500000)
	}
}
