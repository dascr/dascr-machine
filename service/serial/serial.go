package serial

import (
	"bufio"
	"context"
	"fmt"

	"github.com/dascr/dascr-machine/service/config"
	"github.com/dascr/dascr-machine/service/logger"
	"github.com/tarm/serial"
)

// Serial holds the serial connection information
type Serial struct {
	SerialConfig *serial.Config
	Connection   *serial.Port
	Established  bool
	Quit         chan int
	Command      chan string
	Message      chan string
}

// New will return an instantiated Serial object
func New(cmd chan string) *Serial {
	machine := config.Config.Machine
	return &Serial{
		SerialConfig: &serial.Config{
			Name: machine.Serial,
			Baud: 9600,
		},
		Quit:    make(chan int, 1),
		Message: make(chan string),
		Command: cmd,
	}
}

// Start will start the serial loop
func (s *Serial) Start() error {
	ctx, cancel := context.WithCancel(context.Background())

	go s.read(ctx)

	// Infinite read until chan Quit
	for {
		select {
		case c := <-s.Message:
			s.Command <- c
		case <-s.Quit:
			cancel()
			return nil
		}
	}

}

// read will infinitly read on the serial interface
func (s *Serial) read(ctx context.Context) {
	var err error
	s.Established = false
	// Connect via serial
	s.Connection, err = serial.OpenPort(s.SerialConfig)
	if err != nil {
		logger.Errorf("cannot connect to serial: %+v", err)
		config.Config.Machine.Error = err.Error()
		return
	}
	s.Established = true

	s.Write("s,9")
	// Write the Piezo Threshold time to set it at Arduino side
	threshold := fmt.Sprintf("p,%+v", config.Config.Machine.Piezo)
	s.Write(threshold)

	scanner := bufio.NewScanner(s.Connection)

	for {
		select {
		case <-ctx.Done():
			close(s.Message)
			close(s.Quit)
			s.Connection.Close()
			s.Established = false
			return
		default:
			scanner.Scan()

			cmd := scanner.Text()
			if cmd != "" {
				s.Message <- cmd
			}
		}
	}
}

// Reload will terminate the serial loop
func (s *Serial) Reload() {
	logger.Info("Stopping serial loop")
	s.Quit <- 1

	machine := config.Config.Machine

	s.SerialConfig.Name = machine.Serial
	s.Quit = make(chan int, 1)
	s.Message = make(chan string)

	go s.Start()
}

// Write will write to the serial connection
func (s *Serial) Write(input string) {
	b := []byte("<" + input + ">")

	_, err := s.Connection.Write(b)
	if err != nil {
		return
	}
}
