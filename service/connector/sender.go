package connector

import (
	"fmt"
	"log"
	"net/http"
)

func (c *Service) sendToScoreboard(url string) {
	_, err := http.Post(url, "text/plain", nil)
	if err != nil {
		log.Printf("Error sending request %+v to scoreboard: %+v", url, err)
	}
}

func (c *Service) throw(number, modifier int) {
	protocol := "http"
	if c.HTTPS {
		protocol = "https"
	}
	url := fmt.Sprintf("%+v://%+v:%+v/api/game/%+v/throw/%+v/%+v", protocol, c.Host, c.Port, c.Game, number, modifier)
	c.sendToScoreboard(url)
}

func (c *Service) nextPlayer() {
	protocol := "http"
	if c.HTTPS {
		protocol = "https"
	}
	url := fmt.Sprintf("%+v://%+v:%+v/api/game/%+v/nextPlayer", protocol, c.Host, c.Port, c.Game)
	c.sendToScoreboard(url)
}
