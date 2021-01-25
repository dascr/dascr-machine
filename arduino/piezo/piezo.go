package piezo

import "machine"

const threshold = 20

// Piezo holds the pins
type Piezo struct {
	adc *machine.ADC
}

// New will create a new Piezo sensor
func New(input machine.Pin) Piezo {
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

	if val >= threshold {
		return true
	}
	return false
}
