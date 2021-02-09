package connector

import (
	"fmt"

	"github.com/dascr/dascr-machine/service/config"
	"github.com/dascr/dascr-machine/service/game"
	"github.com/dascr/dascr-machine/service/logger"
	"github.com/dascr/dascr-machine/service/sender"
	"github.com/dascr/dascr-machine/service/serial"
	"github.com/dascr/dascr-machine/service/ws"
)

// MachineConnector is the global instance of Connector
var MachineConnector *Connector

// Connector will handle the communication with the arduino
// and sending to the scoreboard
type Connector struct {
	Sender      *sender.Sender
	Websocket   *ws.WS
	Game        *game.Game
	Serial      *serial.Serial
	CommandChan chan string
}

// New will return an instantiated connector object
func New() *Connector {
	cmd := make(chan string)

	serial := serial.New(cmd)
	sender := sender.New(serial)
	game := game.New(cmd, sender, serial)
	ws := ws.New(sender)

	con := &Connector{
		Sender:      sender,
		Game:        game,
		Serial:      serial,
		Websocket:   ws,
		CommandChan: cmd,
	}
	return con

}

// Start will start the connector service
func (c *Connector) Start() error {
	// Debug log the config
	logger.Debugf("Config is: %#v", config.Config)

	go c.Websocket.Start()
	go c.Serial.Start()
	go c.Game.Start()

	logger.Info("Connector service started")

	return nil
}

// ChangePiezoThreshold will write the new piezo threshold
// to the Arduino
func (c *Connector) ChangePiezoThreshold() {
	// Write the Piezo Threshold time to set it at Arduino side
	threshold := fmt.Sprintf("p,%+v", config.Config.Machine.Piezo)
	c.Serial.Write(threshold)
}
