package main

import (
	"machine"

	"github.com/dascr/dascr-machine/arduino/button"
	"github.com/dascr/dascr-machine/arduino/matrix"
	"github.com/dascr/dascr-machine/arduino/piezo"
	"github.com/dascr/dascr-machine/arduino/serial"
)

// constants
const (
	// Ultrasonic
	ultrasonicEchoPin    = machine.D8
	ultrasonicTriggerPin = machine.D9
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
)

func main() {
	machine.InitADC()

	button := button.New(buttonPin, buttonLEDPin)
	button.Configure()

	piezo1 := piezo.New(piezo1Pin)
	piezo2 := piezo.New(piezo2Pin)

	matrix := matrix.New()
	matrix.Configure()

	// Blink 5 times
	button.ButtonLEDBlink(5)

	for {
		if piezo1.ReadPiezo() || piezo2.ReadPiezo() {
			serial.WriteUART("Piezo triggered")
		}
		if button.ReadButton() {
			serial.WriteUART("Button pressed")
		}
		matrixRead := matrix.EvaluateThrow()
		if matrixRead != -1 {
			serial.WriteUART(string(matrixRead))
		}
		serial.ReadUART(button)
	}

	/* example for all sensors
	   machine.InitADC()

	   ultrasonic := ultrasonic.New(ultrasonicTriggerPin, ultrasonicEchoPin)
	   ultrasonic.Configure()

	   button := button.New(buttonPin, buttonLEDPin)
	   button.Configure()

	   piezo1 := piezo.New(piezo1Pin)
	   piezo2 := piezo.New(piezo2Pin)

	   matrix := matrix.New()
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

	/*
		ultrasonic := hcsr04.New(ultrasonicTriggerPin, ultrasonicEchoPin)
		ultrasonic.Configure()

		for {
			distance := ultrasonic.ReadDistance()
			println(distance)
			time.Sleep(time.Millisecond * 100)
		}
	*/

}
