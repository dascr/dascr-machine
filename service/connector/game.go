package connector

import (
	"time"

	"github.com/dascr/dascr-machine/service/logger"
)

// Game will hold minimal state to control machine outputs
type Game struct {
	State string `json:"GameState"`
}

// startGame will be the wrapper to start the main loop as go routine
func (c *Service) startGame() {
	c.mainLoop()
}

// Main loop will handle the machine
func (c *Service) mainLoop() {
	logger.Info("Started main machine loop")
	// Fetch initial state from backend
	c.updateStatus()

	// Loop with select case
	for {
		select {
		case <-c.Quit:
			logger.Info("Canceling main machine loop ...")
			return
		case str := <-c.Command:
			c.processCommand(str)
		default:
		}
	}
}

// processCommand will act on the serial command
// and send to scoreboard depending on command
func (c *Service) processCommand(cmd string) {
	// This is here for DEBUG Purposes and can be removed when tested
	logger.Debugf("Process command %s", cmd)

	// Check if the cmd is in the matrixMap
	if v, ok := matrixMap[cmd]; ok {
		if c.State.State == "THROW" {
			// Send throw
			c.throw(v)
		}
	} else {
		// Else switch over cmd
		switch cmd {
		case "m":
			// missed dart
			if c.State.State == "THROW" {
				// Send miss
				c.throw("0/1")
			}
			break
		case "b":
			// button pressed
			if c.State.State == "WON" {
				c.rematch()
			} else {
				c.nextPlayer()
			}
			break
		case "u":
			// ultrasonic movement
			// As ultrasonic is sending a few more "u" we need to debounce by checking time duration
			if c.State.State != "THROW" && time.Since(c.NextPlayerTime) > c.DebounceTime {
				// Write movement to serial and apply wait delay
				c.Write("s,3")
				// Sleep the Waiting Time from config
				time.Sleep(c.WaitingTime)
				// Then send next player
				c.nextPlayer()
			} else {
				logger.Debug("Not processing it cause of debounce time")
			}
			break

		default:
			break
		}
	}
}
