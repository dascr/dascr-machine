package ws

import (
	"context"
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
	Conn    *websocket.Conn
	Client  *websocket.Dialer
	Header  http.Header
	URL     *url.URL
	Quit    chan int
	Message chan string
	Sender  *sender.Sender
}

// New will return an instantiated websocket connection
func New(sender *sender.Sender) *WS {
	scoreboard := config.Config.Scoreboard

	dialer := &websocket.Dialer{}

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

	if scoreboard.User != "" {
		wsHeaders.Add("Authorization", "Basic "+common.BasicAuth(scoreboard.User, scoreboard.Pass))
	}

	return &WS{
		Conn:    nil,
		Client:  dialer,
		Header:  wsHeaders,
		URL:     u,
		Quit:    make(chan int, 1),
		Message: make(chan string),
		Sender:  sender,
	}
}

// Start will start the websocket connection
func (w *WS) Start() error {
	ctx, cancel := context.WithCancel(context.Background())

	go w.read(ctx)

	// Infinite read until chan Quit
	for {
		select {
		case m := <-w.Message:
			// Handle message
			if m == "update" || m == "redirect" {
				w.Sender.UpdateStatus()
			}
		case <-w.Quit:
			cancel()
			return nil
		}
	}
}

// read will infinitly send messages from ws to channel Message
func (w *WS) read(ctx context.Context) {
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
		return
	}

	logger.Infof("Connected to websocket @ %+v", w.URL.String())
	config.Config.Scoreboard.Error = ""

	logger.Info("Started Websocket listener routine")

	for {
		select {
		case <-ctx.Done():
			close(w.Message)
			close(w.Quit)
			w.Conn.Close()
			return
		default:
			_, message, err := w.Conn.ReadMessage()
			if err != nil {
				logger.Errorf("read:", err)
			}

			w.Message <- string(message)
		}
	}
}

// Reload will reload the websocket connection with new settings
func (w *WS) Reload() {
	logger.Info("Stopping the websocket connection")
	w.Quit <- 1

	scoreboard := config.Config.Scoreboard

	dialer := &websocket.Dialer{}

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

	if scoreboard.User != "" {
		wsHeaders.Add("Authorization", "Basic "+common.BasicAuth(scoreboard.User, scoreboard.Pass))
	}

	w.Conn = nil
	w.Client = dialer
	w.Header = wsHeaders
	w.URL = u
	w.Quit = make(chan int, 1)
	w.Message = make(chan string)

	go w.Start()
}
