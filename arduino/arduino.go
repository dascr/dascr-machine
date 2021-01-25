package main

import (
	"machine"
	"time"

	"github.com/dascr/dascr-machine/arduino/ultrasonic"
)

// constants
const (
	// Ultrasonic
	ultrasonicTriggerPin = machine.D9
	ultrasonicEchoPin    = machine.D8
	// Button
	buttonPin    = machine.D2
	buttonLEDPin = machine.D3
	// Piezo
	piezo1Pin = machine.ADC0
	piezo2Pin = machine.ADC1
)

var (
	matrixOutput = [4]machine.Pin{machine.D22, machine.D24, machine.D49, machine.D47}
	matrixInput  = [16]machine.Pin{machine.D26, machine.D28, machine.D30, machine.D32, machine.D34, machine.D36, machine.D38, machine.D40, machine.D42, machine.D44, machine.D46, machine.D48, machine.D50, machine.D52, machine.D53, machine.D51}
	matrixValues = [4][16]int{
		{212, 112, 209, 109, 214, 114, 211, 111, 208, 108, 000, 312, 309, 314, 311, 308},
		{216, 116, 207, 107, 219, 119, 203, 103, 217, 117, 225, 316, 307, 319, 303, 317},
		{202, 102, 215, 115, 210, 110, 206, 106, 213, 113, 125, 302, 315, 310, 306, 313},
		{204, 104, 218, 118, 201, 101, 220, 120, 205, 105, 000, 304, 318, 301, 320, 305},
	}

	uart = machine.UART0
)

func main() {
	/* echo example
	uart.Write([]byte("Echo console enabled. Type something then press enter:\r\n"))

	input := make([]byte, 64)
	i := 0
	for {
		if uart.Buffered() > 0 {
			data, _ := uart.ReadByte()
			switch data {
			case 13:
				// return key
				uart.Write([]byte("\r\n"))
				uart.Write([]byte("You typed: "))
				uart.Write(input[:i])
				uart.Write([]byte("\r\n"))
				i = 0
			default:
				// just echo the character
				uart.WriteByte(data)
				input[i] = data
				i++
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
	*/

	/* example for all sensors
	machine.InitADC()

	ultrasonic := ultrasonic.New(ultrasonicTriggerPin, ultrasonicEchoPin)
	ultrasonic.Configure()

	button := button.New(buttonPin, buttonLEDPin)
	button.Configure()

	piezo1 := piezo.New(piezo1Pin)
	piezo2 := piezo.New(piezo2Pin)

	matrix := matrix.New(matrixOutput, matrixInput, matrixValues)
	matrix.Configure()

	for {
		println("Distance: ", ultrasonic.ReadDistance())
		if piezo1.ReadPiezo() || piezo2.ReadPiezo() {
			println("Piezo triggered")
		}
		if button.ReadButton() {
			println("Button pressed")
		}
		matrixRead := matrix.EvaluateThrow()
		if matrixRead != -1 {
			println("Matrix: ", matrixRead)
		}
		time.Sleep(time.Millisecond * 1000)
	}
	*/

	ultrasonic := ultrasonic.New(ultrasonicTriggerPin, ultrasonicEchoPin)
	ultrasonic.Configure()

	for {
		println("Distance: ", ultrasonic.ReadDistance())
		time.Sleep(time.Millisecond * 100)
	}

}
