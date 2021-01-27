package connector

import (
	"log"
	"time"

	"github.com/dascr/dascr-machine/service/config"
	"github.com/tarm/serial"
)

// Serv is the global connector service
var Serv Service

// Service will handle the communication with the arduino
// and sending to the scoreboard
type Service struct {
	WaitingTime int
	Host        string
	Port        string
	Config      *serial.Config
	Conn        *serial.Port
}

// Start will start the connector service
func (c *Service) Start() error {
	// Read config
	c.WaitingTime = config.Config.Machine.WaitingTime
	c.Host = config.Config.Scoreboard.Host
	c.Port = config.Config.Scoreboard.Port

	c.Config.Name = config.Config.Machine.Serial

	// Connect via serial
	port, err := serial.OpenPort(c.Config)
	if err != nil {
		return err
	}

	c.Conn = port

	// Write 7 to blink 7 times to serial
	c.Conn.Write([]byte("7\r"))

	log.Println("Connector service started")

	return nil
}

// Stop will stop the connector service
func (c *Service) Stop() {
	log.Println("Connector service stopped")
	c.Conn.Close()
}

// Restart will first stop then start the connector service again
func (c *Service) Restart() error {
	log.Println("Restarting connector service")
	c.Stop()
	time.Sleep(time.Second * 2)
	err := c.Start()
	if err != nil {
		return err
	}
	return nil
}

// Write will write to the serial connection
func (c *Service) Write(input string) {
	b := []byte(input + "\r")

	_, err := c.Conn.Write(b)
	if err != nil {
		log.Printf("Error writing to serial connection: %+v", err)
		return
	}
}

// Read will read from the serial connection
func (c *Service) Read() string {
	buf := make([]byte, 128)
	n, _ := c.Conn.Read(buf)
	log.Printf("Buf is %+v", buf)
	output := string(buf[:n])
	log.Printf("Read output in function is: %+v", output)

	return output
}
