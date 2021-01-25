package main

import (
	"machine"
	"time"
)

// LEDButton holds the pins
type LEDButton struct {
	button machine.Pin
	led    machine.Pin
}

// NewLEDButton creates a new led button
func NewLEDButton(button, led machine.Pin) LEDButton {
	return LEDButton{
		button: button,
		led:    led,
	}
}

// Configure will setup the button
func (lb *LEDButton) Configure() {
	lb.button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	lb.led.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

// ReadButton will return bool of button state
func (lb *LEDButton) ReadButton() bool {
	return lb.button.Get()
}

// ButtonLEDOn will turn the buttons LED on
func (lb *LEDButton) ButtonLEDOn() {
	lb.led.Low()
}

// ButtonLEDOff will turn the buttons LED off
func (lb *LEDButton) ButtonLEDOff() {
	lb.led.High()
}

// Ultrasonic holds the pins
type Ultrasonic struct {
	trigger machine.Pin
	echo    machine.Pin
}

// NewUltrasonic creates a new ultrasonic sensor
func NewUltrasonic(trigger, echo machine.Pin) Ultrasonic {
	return Ultrasonic{
		trigger: trigger,
		echo:    echo,
	}
}

// Configure will setup the ultrasonic sensor
func (u *Ultrasonic) Configure() {
	u.trigger.Configure(machine.PinConfig{Mode: machine.PinOutput})
	u.echo.Configure(machine.PinConfig{Mode: machine.PinInput})
}

// ReadDistance reads the distance and returns it in mm
func (u *Ultrasonic) ReadDistance() int32 {
	pulse := u.ReadPulse()

	return (pulse * 1715) / 10000 // mm
}

// ReadPulse will read the pulse of the ultrasonic sensor
func (u *Ultrasonic) ReadPulse() int32 {
	t := time.Now()
	u.trigger.Low()
	time.Sleep(2 * time.Microsecond)
	u.trigger.High()
	time.Sleep(10 * time.Microsecond)
	u.trigger.Low()
	i := uint8(0)
	for {
		if u.echo.Get() {
			t = time.Now()
			break
		}
		i++
		if i > 10 {
			if time.Since(t).Microseconds() > timeout {
				return 0
			}
			i = 0
		}
	}
	i = 0
	for {
		if !u.echo.Get() {
			return int32(time.Since(t).Microseconds())
		}
		i++
		if i > 10 {
			if time.Since(t).Microseconds() > timeout {
				return 0
			}
			i = 0
		}
	}
}

// Piezo holds the pins
type Piezo struct {
	adc *machine.ADC
}

// NewPiezo will create a new Piezo sensor
func NewPiezo(input machine.Pin) Piezo {
	adc := machine.ADC{input}
	return Piezo{
		adc: &adc,
	}
}

// ReadPiezo will read the piezo sensor and compare it to the treshold
// and finally return a bool
func (p *Piezo) ReadPiezo() bool {
	val := p.adc.Get()
	val = p.adc.Get()

	if val >= pThreshold {
		return true
	}
	return false
}

// Matrix holds the pins
type Matrix struct {
	outputs [4]machine.Pin
	inputs  [16]machine.Pin
	values  [4][16]int
}

// NewMatrix will create a new matrix sensor
func NewMatrix(outputs [4]machine.Pin, inputs [16]machine.Pin, values [4][16]int) Matrix {
	return Matrix{
		outputs: outputs,
		inputs:  inputs,
		values:  values,
	}
}

// Configure will setup all the pins you need for matrix evaluation
func (m *Matrix) Configure() {
	// Set outputs to output
	for _, pin := range m.outputs {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}
	// Set inputs to input pullup
	for _, pin := range m.inputs {
		pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}
}

// EvaluateThrow will return the value of the throw if any
func (m *Matrix) EvaluateThrow() int {
	var eval [4][16]bool
	for i := 0; i < 4; i++ {
		m.outputs[0].High()
		m.outputs[1].High()
		m.outputs[2].High()
		m.outputs[3].High()
		m.outputs[i].Low()

		for j := 0; j < 16; j++ {
			eval[i][j] = m.inputs[j].Get()
			if !eval[i][j] {
				return m.values[i][j]
			}
		}
	}
	return -1
}

// constants
const (
	// Ultrasonic
	timeout              = 23324 // 4m
	ultrasonicTriggerPin = machine.D9
	ultrasonicEchoPin    = machine.D8
	// Button
	buttonPin    = machine.D2
	buttonLEDPin = machine.D3
	// Piezo
	pThreshold = 20
	piezo1Pin  = machine.ADC0
	piezo2Pin  = machine.ADC1
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
)

func main() {

	ultrasonic := NewUltrasonic(ultrasonicTriggerPin, ultrasonicEchoPin)
	ultrasonic.Configure()

	for {
		println(ultrasonic.ReadDistance())
		time.Sleep(time.Millisecond * 100)
	}

	/* Button test
	ledbutton := NewLEDButton(buttonPin, buttonLEDPin)
	ledbutton.Configure()

	for {
		if ledbutton.ReadButton() {
			ledbutton.ButtonLEDOn()
		} else {
			ledbutton.ButtonLEDOff()
		}
	}
	*/

	/* Piezo test
	machine.InitADC()

	ledbutton := NewLEDButton(buttonPin, buttonLEDPin)
	ledbutton.Configure()
	piezo1 := NewPiezo(piezo1Pin)
	piezo2 := NewPiezo(piezo2Pin)

	for {
		if piezo1.ReadPiezo() || piezo2.ReadPiezo() {
			ledbutton.ButtonLEDOn()

		} else {
			ledbutton.ButtonLEDOff()
		}
		time.Sleep(100 * time.Millisecond)
	}
	*/

	/* Matrix test
	ledbutton := NewLEDButton(buttonPin, buttonLEDPin)
	ledbutton.Configure()

	matrix := NewMatrix(matrixOutput, matrixInput, matrixValues)
	matrix.Configure()

	for {
		matrixVal := matrix.EvaluateThrow()
		if matrixVal != -1 {
			println(matrixVal)
		}

		time.Sleep(time.Millisecond * 100)
	}
	*/
}
