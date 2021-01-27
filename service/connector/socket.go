package connector

import (
	"log"
	"time"
)

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
		log.Printf("Got websocket message: %+v", string(message))
		if string(message) == "update" {
			time.Sleep(time.Second * 1)
			c.updateStatus()
			log.Println("Status updated")
		}
	}
}
