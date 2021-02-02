package connector

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/dascr/dascr-machine/service/config"
	"github.com/dascr/dascr-machine/service/logger"
	"github.com/gorilla/websocket"
	"github.com/tarm/serial"
)

// Serv is the global connector service
var Serv Service

// Service will handle the communication with the arduino
// and sending to the scoreboard
type Service struct {
	WaitingTime     time.Duration
	DebounceTime    time.Duration
	NextPlayerTime  time.Time
	Piezo           int
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
	HTTPClient      *http.Client
	CookieJar       *cookiejar.Jar
	WebsocketClient *websocket.Dialer
}

// Start will start the connector service
func (c *Service) Start() error {
	// Debug log the config
	logger.Debugf("Config is: %#v", config.Config)
	if !c.Running {
		var err error
		var resp *http.Response
		// Read config
		c.WaitingTime = time.Duration(config.Config.Machine.WaitingTime * int(time.Second))
		c.DebounceTime = time.Duration(c.WaitingTime + 2*time.Second)
		c.NextPlayerTime = time.Now()
		c.Piezo = config.Config.Machine.Piezo
		c.Config.Name = config.Config.Machine.Serial

		c.HTTPS = config.Config.Scoreboard.HTTPS
		c.Host = config.Config.Scoreboard.Host
		c.Port = config.Config.Scoreboard.Port
		c.Game = config.Config.Scoreboard.Game
		c.User = config.Config.Scoreboard.User
		c.Pass = config.Config.Scoreboard.Pass

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
				logger.Errorf("handshake failed with status %d", resp.StatusCode)
			}
			logger.Errorf("cannot connect to websocket: %+v", err)
			config.Config.Scoreboard.Error = err.Error()
			return err
		}
		logger.Infof("Connected to websocket @ %+v", u.String())
		config.Config.Scoreboard.Error = ""

		// Connect via serial
		c.Conn, err = serial.OpenPort(c.Config)
		if err != nil {
			logger.Errorf("cannot connect to serial: %+v", err)
			config.Config.Machine.Error = err.Error()
			return err
		}

		logger.Info("Serial connection initiated")
		config.Config.Machine.Error = ""

		// Write 9 to arduino indicating we are connected
		// Button should blink 7 times
		c.Write("s,9")
		c.Running = true
		// Write the Piezo Threshold time to set it at Arduino side
		threshold := fmt.Sprintf("p,%+v", c.Piezo)
		c.Write(threshold)

		logger.Info("Connector service started")

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
// func (c *Service) Stop(ctx context.Context) {
func (c *Service) Stop() {
	// Only stop if running
	if c.Running {
		// Terminate go routines
		c.Quit <- 0
		// Shutdown serial connection
		c.Conn.Close()

		c.Running = false
		logger.Info("Connector service stopped")
	}
}

// Restart will first stop then start the connector service again
func (c *Service) Restart() error {
	// Only stop if running
	if c.Running {
		logger.Info("Restarting connector service")
		// ctx, cancle := context.WithTimeout(context.Background(), 15)
		// defer cancle()
		// c.Stop(ctx)
		c.Stop()
		time.Sleep(time.Second * 3)
	}

	err := c.Start()
	if err != nil {
		return err
	}
	return nil
}
