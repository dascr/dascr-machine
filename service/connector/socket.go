package connector

import (
	"log"
)

func (c *Service) listenToWebsocket() {
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
		if string(message) == "update" {
			c.updateStatus()
		}
	}
}
