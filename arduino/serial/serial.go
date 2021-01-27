package serial

import (
	"machine"
	"time"

	"github.com/dascr/dascr-machine/arduino/button"
)

var uart = machine.UART0

// ReadUART will read if there is something send to the board
// and then dispatch the action
func ReadUART(button button.LEDButton) {
	input := make([]byte, 64)
	i := 0
	for {
		if uart.Buffered() > 0 {
			data, _ := uart.ReadByte()
			switch data {
			case 13:
				// return key aka \r
				cmd := string(input[:i])
				if cmd == "1" {
					button.ButtonLEDOn()
				} else if cmd == "2" {
					button.ButtonLEDOff()
				} else if cmd == "3" {
					button.ButtonLEDBlink(3)
				} else if cmd == "7" {
					button.ButtonLEDBlink(7)
				}
				i = 0
			default:
				input[i] = data
				i++
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// WriteUART writes back to uart
func WriteUART(output string) {
	uart.Write([]byte(output))
}
