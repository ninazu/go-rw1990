package main

import "C"
import (
	"fmt"
	"github.com/janne/bcm2835"
	"time"
)

const PIN = bcm2835.Pin11

const (
	BCM2835_GPIO_FSEL_INPT = 0
	BCM2835_GPIO_FSEL_OUTP = 1
	READ_OPCODE            = 0x33
	WRITE_OPCODE           = 0xD5
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

func writeCode(data []byte) bool {
	sendR()

	if !waitP() {
		return false
	}

	send(WRITE_OPCODE)

	for bit := 0; bit < 8; bit++ {
		write(data[bit])
	}

	delay(16000)

	sendR()
	waitP()

	return true
}

func sendR() {
	down()
	delay(1024)
	up()
	delay(1)
}

func waitP() bool {
	count := 10000

	for bcm2835.GpioLev(PIN) == 0 {
		count--

		if count < 0 {
			return false
		}
	}

	for bcm2835.GpioLev(PIN) != 0 {
		delay(1)
	}

	delay(1)

	return true
}

func receive() byte {
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

func readCode() (bool, []byte) {
	sendR()

	if !waitP() {
		return false, nil
	}

	send(READ_OPCODE)
	var data = make([]byte, 8)

	for bit := 0; bit < 8; bit++ {
		data[bit] = receive()
	}

	return true, data
}

func readButton() []byte {
	var data []byte

	for {
		if result, data := readCode(); result {
			fmt.Printf("Read: %x\n", data)
			break
		} else {
			fmt.Printf(".\n")
		}

		delay(500000)
	}

	return data
}

func writeButton(data []byte) bool {
	if writeCode(data) {
		fmt.Printf("Write: %x\n", data)
	} else {
		fmt.Printf("Writing failed \n")

		return false
	}

	return true
}

func copyButton() {
	data := readButton()

	if !writeButton(data) {
		return
	}

	fmt.Printf("Verifying\n")
	readButton()
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

	readButton()
}
