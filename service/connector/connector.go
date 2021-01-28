package connector

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
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
	WaitingTime     int
	HTTPS           bool
	Host            string
	Port            string
	Game            string
	User            string
	Pass            string
	Config          *serial.Config
	Conn            *serial.Port
	WebsocketConn   *websocket.Conn
	Running         bool
	Quit            chan int
	Command         chan string
	State           Game
	Scanner         *bufio.Scanner
	HTTPClient      *http.Client
	CookieJar       *cookiejar.Jar
	WebsocketClient *websocket.Dialer
}

// Start will start the connector service
func (c *Service) Start() error {
	if !c.Running {
		var err error
		var resp *http.Response
		// Read config
		c.WaitingTime = config.Config.Machine.WaitingTime
		c.Config.Name = config.Config.Machine.Serial

		c.HTTPS = config.Config.Scoreboard.HTTPS
		c.Host = config.Config.Scoreboard.Host
		c.Port = config.Config.Scoreboard.Port
		c.Game = config.Config.Scoreboard.Game

		c.CookieJar = &cookiejar.Jar{}

		c.HTTPClient = &http.Client{
			Jar:           c.CookieJar,
			CheckRedirect: c.redirectPolicyFunc,
		}

		c.WebsocketClient = &websocket.Dialer{
			Jar: c.CookieJar,
		}

		// Other vars
		scheme := "ws"
		if c.HTTPS {
			scheme = "wss"
		}

		host := fmt.Sprintf("%+v", c.Host)
		if c.Port != "80" && c.Port != "443" {
			host = fmt.Sprintf("%+v:%+v", c.Host, c.Port)
		}
		path := fmt.Sprintf("/ws/%+v", c.Game)

		u := &url.URL{Scheme: scheme, Host: host, Path: path}
		log.Printf("DEBUG: ws url is: %+v", u.String())

		// create channels
		c.Quit = make(chan int)
		c.Command = make(chan string)

		// Init state
		c.State = Game{}

		wsHeaders := &http.Header{}
		// Connect Websocket
		if c.User != "" {
			wsHeaders.Add("Authorization", "Basic "+c.basicAuth())
		}

		c.WebsocketConn, resp, err = c.WebsocketClient.Dial(u.String(), *wsHeaders)
		if err != nil {
			if err == websocket.ErrBadHandshake {
				log.Printf("handshake failed with status %d", resp.StatusCode)
			}
			log.Printf("cannot connect to websocket: %+v", err)
			config.Config.Scoreboard.Error = err.Error()
			return err
		}
		log.Printf("Connected to websocket @ %+v", u.String())
		config.Config.Scoreboard.Error = ""

		// Connect via serial
		c.Conn, err = serial.OpenPort(c.Config)
		if err != nil {
			log.Printf("cannot connect to serial: %+v", err)
			config.Config.Machine.Error = err.Error()
			return err
		}
		// Assign scanner
		c.Scanner = bufio.NewScanner(c.Conn)

		log.Println("Serial connection initiated")
		config.Config.Machine.Error = ""

		// Write 7 to blink 7 times to serial
		c.Conn.Write([]byte("7\n"))
		c.Running = true

		log.Println("Connector service started")

		// Start websocket
		go c.startWebsocket()

		// Start serial reader
		go c.startSerial()

		// start main loop
		go c.startGame()
	}

	return nil
}

// Stop will stop the connector service
func (c *Service) Stop(ctx context.Context) {
	// Only stop if running
	if c.Running {
		// Terminate go routines
		c.Quit <- 0
		// Shutdown serial connection
		c.Conn.Close()

		c.Running = false
		log.Println("Connector service stopped")
	}
}

// Restart will first stop then start the connector service again
func (c *Service) Restart() error {
	// Only stop if running
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
