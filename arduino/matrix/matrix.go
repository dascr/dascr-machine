package matrix

import "machine"

var (
	matrixValues = [4][16]int{
		{212, 112, 209, 109, 214, 114, 211, 111, 208, 108, 000, 312, 309, 314, 311, 308},
		{216, 116, 207, 107, 219, 119, 203, 103, 217, 117, 225, 316, 307, 319, 303, 317},
		{202, 102, 215, 115, 210, 110, 206, 106, 213, 113, 125, 302, 315, 310, 306, 313},
		{204, 104, 218, 118, 201, 101, 220, 120, 205, 105, 000, 304, 318, 301, 320, 305},
	}

	matrixOutput = [4]machine.Pin{machine.D22, machine.D24, machine.D49, machine.D47}

	matrixInput = [16]machine.Pin{machine.D26, machine.D28, machine.D30, machine.D32, machine.D34, machine.D36, machine.D38, machine.D40, machine.D42, machine.D44, machine.D46, machine.D48, machine.D50, machine.D52, machine.D53, machine.D51}
)

// Matrix holds the pins
type Matrix struct {
	outputs [4]machine.Pin
	inputs  [16]machine.Pin
	values  [4][16]int
}

// New will create a new matrix sensor
func New() Matrix {
	return Matrix{
		outputs: matrixOutput,
		inputs:  matrixInput,
		values:  matrixValues,
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
func (m *Matrix) EvaluateThrow() string {
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
				return string(m.values[i][j])
			}
		}
	}
	return "-1"
}
