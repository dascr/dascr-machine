package matrix

import "machine"

// Matrix holds the pins
type Matrix struct {
	outputs [4]machine.Pin
	inputs  [16]machine.Pin
	values  [4][16]int
}

// New will create a new matrix sensor
func New(outputs [4]machine.Pin, inputs [16]machine.Pin, values [4][16]int) Matrix {
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
