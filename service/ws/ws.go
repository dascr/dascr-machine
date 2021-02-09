package ws

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/dascr/dascr-machine/service/common"
	"github.com/dascr/dascr-machine/service/config"
	"github.com/dascr/dascr-machine/service/logger"
	"github.com/dascr/dascr-machine/service/sender"
	"github.com/gorilla/websocket"
)

// WS will hold the websocket information
type WS struct {
	Conn   *websocket.Conn
	Client *websocket.Dialer
	Header http.Header
	URL    *url.URL
	Quit   chan int
	Sender *sender.Sender
}

// New will return an instantiated websocket connection
func New(sender *sender.Sender) *WS {
	scoreboard := config.Config.Scoreboard

	dialer := &websocket.Dialer{}

	// Other vars
	scheme := "ws"
	if scoreboard.HTTPS {
		scheme = "wss"
	}

	host := fmt.Sprintf("%+v", scoreboard.Host)
	if scoreboard.Port != "80" && scoreboard.Port != "443" {
		host = fmt.Sprintf("%+v:%+v", scoreboard.Host, scoreboard.Port)
	}
	path := fmt.Sprintf("/ws/%+v", scoreboard.Game)

	u := &url.URL{Scheme: scheme, Host: host, Path: path}

	wsHeaders := http.Header{}
	// Connect Websocket
	if scoreboard.User != "" {
		wsHeaders.Add("Authorization", "Basic "+common.BasicAuth(scoreboard.User, scoreboard.Pass))
	}

	return &WS{
		Conn:   nil,
		Client: dialer,
		Header: wsHeaders,
		URL:    u,
		Quit:   make(chan int),
		Sender: sender,
	}
}

// Start will start the websocket connection
func (w *WS) Start() error {
	var resp *http.Response
	var err error

	// Connect
	w.Conn, resp, err = w.Client.Dial(w.URL.String(), w.Header)
	if err != nil {
		if err == websocket.ErrBadHandshake {
			logger.Errorf("handshake failed with status %d", resp.StatusCode)
		}
		logger.Errorf("cannot connect to websocket: %+v", err)
		config.Config.Scoreboard.Error = err.Error()
		return err
	}

	logger.Infof("Connected to websocket @ %+v", w.URL.String())
	config.Config.Scoreboard.Error = ""

	logger.Info("Started Websocket listener routine")

	// Infinite read until chan Quit
	for {
		select {
		case <-w.Quit:
			return nil
		default:
		}
		_, message, err := w.Conn.ReadMessage()
		if err != nil {
			logger.Errorf("read:", err)
			return err
		}

		// Handle message
		if string(message) == "update" || string(message) == "redirect" {
			w.Sender.UpdateStatus()
		}
	}
}
