package connector

import (
	"log"
)

// startWebsocket is a wrapper to start the websocket client as go routine
func (c *Service) startWebsocket() {
	c.listenToWebsocket()
}

// listentoWebsocket will listen for updates
// and then update the status
func (c *Service) listenToWebsocket() {
	log.Println("Started Websocket listener routine")
	for {
		select {
		case <-c.Quit:
			return
		default:
		}
		_, message, err := c.WebsocketConn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		if string(message) == "update" || string(message) == "redirect" {
			c.updateStatus()
		}
	}
}
