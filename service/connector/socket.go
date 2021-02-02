package connector

import "github.com/dascr/dascr-machine/service/logger"

// startWebsocket is a wrapper to start the websocket client as go routine
func (c *Service) startWebsocket() {
	c.listenToWebsocket()
}

// listentoWebsocket will listen for updates
// and then update the status
func (c *Service) listenToWebsocket() {
	logger.Info("Started Websocket listener routine")
	for {
		select {
		case <-c.Quit:
			return
		default:
		}
		_, message, err := c.WebsocketConn.ReadMessage()
		if err != nil {
			logger.Errorf("read:", err)
			return
		}

		if string(message) == "update" || string(message) == "redirect" {
			c.updateStatus()
			if c.State.State == "" {
				c.buttonOn()
			}
		}
	}
}
