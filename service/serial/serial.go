package serial

import (
	"bufio"
	"fmt"

	"github.com/dascr/dascr-machine/service/config"
	"github.com/dascr/dascr-machine/service/logger"
	"github.com/tarm/serial"
)

// Serial holds the serial connection information
type Serial struct {
	SerialConfig *serial.Config
	Connection   *serial.Port
	Quit         chan int
	Command      chan string
}

// New will return an instantiated Serial object
func New(cmd chan string) *Serial {
	machine := config.Config.Machine
	return &Serial{
		SerialConfig: &serial.Config{
			Name: machine.Serial,
			Baud: 9600,
		},
		Quit:    make(chan int),
		Command: cmd,
	}
}

// Start will start the serial loop
func (s *Serial) Start() error {
	var err error
	// Connect via serial
	s.Connection, err = serial.OpenPort(s.SerialConfig)
	if err != nil {
		logger.Errorf("cannot connect to serial: %+v", err)
		config.Config.Machine.Error = err.Error()
		return err
	}

	s.Write("s,9")
	// Write the Piezo Threshold time to set it at Arduino side
	threshold := fmt.Sprintf("p,%+v", config.Config.Machine.Piezo)
	s.Write(threshold)

	scanner := bufio.NewScanner(s.Connection)

	for scanner.Scan() {
		select {
		case <-s.Quit:
			s.Connection.Close()
			return nil
		default:
		}

		cmd := scanner.Text()
		if cmd != "" {
			s.Command <- cmd
		}
	}
	return nil
}

// Reload will terminate the serial loop
func (s *Serial) Reload() {
	logger.Error("Still need to implement serial reload")
}

// Write will write to the serial connection
func (s *Serial) Write(input string) {
	b := []byte("<" + input + ">")

	_, err := s.Connection.Write(b)
	if err != nil {
		return
	}
}
