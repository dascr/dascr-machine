package connector

import (
	"log"
	"time"
)

// Game will hold minimal state to control machine outputs
type Game struct {
	State           string `json:"GameState"`
	RetrievingDarts bool
}

// startGame will be the wrapper to start the main loop as go routine
func (c *Service) startGame() {
	c.mainLoop()
}

// Main loop will handle the machine
func (c *Service) mainLoop() {
	log.Println("Started main machine loop")
	// Fetch initial state from backend
	c.updateStatus()

	// Loop with select case
	for {
		select {
		case <-c.Quit:
			log.Println("Canceling main machine loop ...")
			return
		case str := <-c.Command:
			c.processCommand(str)
		default:
		}

		// Default is to update button state
		// c.stateMachine()
	}
}

/*
// stateMachine will handle LED of button
func (c *Service) stateMachine() {
	switch c.State.State {
	case "WON":
		c.buttonBlink1()
	case "NEXTPLAYER":
		c.buttonOn()
	case "THROW":
		c.buttonOff()
	case "BUST":
		c.buttonOn()
	case "BUSTCONDITION":
		c.buttonOn()
	case "BUSTNOCHECKOUT":
		c.buttonOn()
	default:
		break
	}
}
*/

// processCommand will act on the serial command
// and send to scoreboard depending on command
func (c *Service) processCommand(cmd string) {
	// This is here for DEBUG Purposes and can be removed when tested
	log.Printf("DEBUG: Process command %s", cmd)

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
		case "b":
			// button pressed
			if c.State.State == "WON" {
				c.rematch()
			} else {
				c.nextPlayer()
			}
		case "u":
			// ultrasonic movement
			if c.State.State != "THROW" {
				// Write movement to serial and apply wait delay
				c.Write("3")
				// Sleep the Waiting Time from config
				time.Sleep(time.Second * time.Duration(c.WaitingTime))
				// Then send next player
				c.nextPlayer()
			}

		default:
			break
		}
	}
}
