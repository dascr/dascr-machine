package game

import (
	"time"

	"github.com/dascr/dascr-machine/service/config"
	"github.com/dascr/dascr-machine/service/logger"
	"github.com/dascr/dascr-machine/service/sender"
	"github.com/dascr/dascr-machine/service/serial"
	"github.com/dascr/dascr-machine/service/state"
)

// Game will hold minimal state to control machine outputs
type Game struct {
	WaitingTime    time.Duration
	DebounceTime   time.Duration
	NextPlayerTime time.Time
	Quit           chan int
	Command        chan string
	Sender         *sender.Sender
	Serial         *serial.Serial
}

// New will return an instantiated Game
func New(cmd chan string, sender *sender.Sender, serial *serial.Serial) *Game {
	machine := config.Config.Machine
	w := time.Second * time.Duration(machine.WaitingTime)
	d := w + 2*time.Second

	game := &Game{
		WaitingTime:    w,
		DebounceTime:   d,
		NextPlayerTime: time.Now(),
		Quit:           make(chan int),
		Command:        cmd,
		Sender:         sender,
		Serial:         serial,
	}
	return game
}

// Start will be the wrapper to start the main loop as go routine
func (g *Game) Start() {
	logger.Info("Started main machine loop")
	// Fetch initial state from backend
	g.Sender.UpdateStatus()

	// Loop with select case
	for {
		select {
		case <-g.Quit:
			logger.Info("Canceling main machine loop ...")
			return
		case str := <-g.Command:
			g.processCommand(str)
		default:
		}
	}
}

// processCommand will act on the serial command
// and send to scoreboard depending on command
func (g *Game) processCommand(cmd string) {
	// This is here for DEBUG Purposes and can be removed when tested
	logger.Debugf("Process command %s", cmd)

	// Check if the cmd is in the matrixMap
	if v, ok := matrixMap[cmd]; ok {
		if state.GameState.GameState == "THROW" {
			// Send throw
			g.Sender.Throw(v)
		}
	} else {
		// Else switch over cmd
		switch cmd {
		case "m":
			// missed dart
			if state.GameState.GameState == "THROW" {
				// Send miss
				g.Sender.Throw("0/1")
			}
			break
		case "b":
			// button pressed
			if state.GameState.GameState == "WON" {
				g.Sender.Rematch()
			} else {
				g.Sender.NextPlayer()
			}
			break
		case "u":
			// ultrasonic movement
			// As ultrasonic is sending a few more "u" we need to debounce by checking time duration
			if state.GameState.GameState != "THROW" && time.Since(g.NextPlayerTime) > time.Second*9 {
				// Write movement to serial and apply wait delay
				g.Serial.Write("s,3")
				// Sleep the Waiting Time from config
				time.Sleep(g.WaitingTime)
				// Write 4 to serial to reset ultrasonic loop at Arduino
				g.Serial.Write("s,4")
				// Then send next player
				g.Sender.NextPlayer()
				// Setting the time after nextPlayer to debounce multiple "u"'s
				g.NextPlayerTime = time.Now()
			} else {
				logger.Debug("Not processing it cause of debounce time")
			}
			break

		default:
			break
		}
	}

}
