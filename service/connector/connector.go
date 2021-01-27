package connector

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/dascr/dascr-machine/service/config"
	"github.com/gorilla/websocket"
	"github.com/tarm/serial"
)

// Serv is the global connector service
var Serv Service

// Service will handle the communication with the arduino
// and sending to the scoreboard
type Service struct {
	WaitingTime   int
	HTTPS         bool
	Host          string
	Port          string
	Game          string
	Config        *serial.Config
	Conn          *serial.Port
	WebsocketConn *websocket.Conn
	Running       bool
	Quit          chan int
	State         Game
}

// Start will start the connector service
func (c *Service) Start() error {
	if !c.Running {
		var err error
		// Read config
		c.WaitingTime = config.Config.Machine.WaitingTime
		c.Config.Name = config.Config.Machine.Serial

		c.HTTPS = config.Config.Scoreboard.HTTPS
		c.Host = config.Config.Scoreboard.Host
		c.Port = config.Config.Scoreboard.Port
		c.Game = config.Config.Scoreboard.Game

		// Other vars
		scheme := "ws"
		if c.HTTPS {
			scheme = "wss"
		}
		host := fmt.Sprintf("%+v:%+v", c.Host, c.Port)
		path := fmt.Sprintf("/ws/%+v", c.Game)

		u := url.URL{Scheme: scheme, Host: host, Path: path}

		// channel
		c.Quit = make(chan int)

		// Init state
		c.State = Game{}

		// Connect Websocket
		c.WebsocketConn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			return err
		}

		// Connect via serial
		c.Conn, err = serial.OpenPort(c.Config)
		if err != nil {
			return err
		}

		// Write 7 to blink 7 times to serial
		c.Conn.Write([]byte("7\n"))
		c.Running = true

		log.Println("Connector service started")

		// Start websocket
		go func() {
			c.listenToWebsocket()
		}()

		// start main loop
		go func() {
			c.mainLoop()
		}()
	}

	return nil
}

// Stop will stop the connector service
func (c *Service) Stop(ctx context.Context) {
	if c.Running {
		c.Quit <- 0
		c.Conn.Close()
		c.Running = false
		log.Println("Connector service stopped")
	}
}

// Restart will first stop then start the connector service again
func (c *Service) Restart() error {
	if c.Running {
		log.Println("Restarting connector service")
		ctx, cancle := context.WithTimeout(context.Background(), 15)
		defer cancle()
		c.Stop(ctx)
		time.Sleep(time.Second * 2)
	}
	err := c.Start()
	if err != nil {
		return err
	}
	return nil
}

// Write will write to the serial connection
func (c *Service) Write(input string) {
	b := []byte(input + "\n")

	_, err := c.Conn.Write(b)
	if err != nil {
		log.Printf("Error writing to serial connection: %+v", err)
		return
	}
}

// Read will read from the serial connection
func (c *Service) Read() string {
	var cmd string

	scanner := bufio.NewScanner(c.Conn)
	for scanner.Scan() {
		cmd = scanner.Text()
	}

	return cmd
}

func (c *Service) mainLoop() {
	// Fetch initial state from backend
	c.updateStatus()

	for {
		select {
		case <-c.Quit:
			log.Println("Canceling main machine loop ...")
			return
		default:
		}

		// State machine
		switch c.State.State {
		case "WON":
			c.buttonOn()
		case "NEXTPLAYER":
			c.buttonOn()
		case "THROW":
			c.buttonOff()
		case "BUST":
			c.buttonOn()
		default:
			break
		}

		// This is the actual machine loop
		cmd := c.Read()

		if cmd != "" {
			log.Printf("Machine said: %+v", cmd)

			// Check if the cmd is in the matrixMap
			if v, ok := matrixMap[cmd]; ok {
				log.Println("Hitting throw")
				// Send throw
				c.throw(v)
			} else {
				// Else switch over cmd
				switch cmd {
				case "m":
					log.Println("Hitting missed")
					// missed dart
					c.throw("0/1")
				case "b":
					log.Println("Hitting button")
					// button pressed
					if c.State.State == "WON" {
						c.rematch()
					} else {
						c.nextPlayer()
					}
				case "u":
					log.Println("Hitting ultrasonic")
					// ultrasonic movement
					// Write movement to serial and apply wait delay
					c.Write("3")
					time.Sleep(time.Second * time.Duration(c.WaitingTime))
					// Then send next player
					c.nextPlayer()
					// Turn button off
					c.buttonOff()

				default:
					break
				}
			}
		}
	}
}
