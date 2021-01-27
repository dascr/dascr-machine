package serial

import (
	"machine"
	"time"
)

var uart = machine.UART0

// ReadUART will read if there is something send to the board
// and then dispatch the action
func ReadUART() string {
	input := make([]byte, 64)
	i := 0
	for {
		if uart.Buffered() > 0 {
			data, _ := uart.ReadByte()
			switch data {
			case 13:
				// return key aka \r
				cmd := string(input[:i])
				return cmd
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
	write := []byte(output)
	uart.Write(write)
}
