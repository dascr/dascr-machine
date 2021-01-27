package connector

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Game will hold minimal state to control machine outputs
type Game struct {
	State   string `json:"GameState"`
	Message string `json:"Message"`
}

func (c *Service) sendToScoreboard(url, method string) (*http.Response, error) {
	var resp *http.Response
	var err error

	protocol := "http"
	if c.HTTPS {
		protocol = "https"
	}

	target := fmt.Sprintf("%+v://%+v:%+v/api/game/%+v/%+v", protocol, c.Host, c.Port, c.Game, url)

	switch method {
	case "get":
		resp, err = http.Get(target)
		if err != nil {
			return nil, err
		}
	case "post":
		resp, err = http.Post(target, "text/plain", nil)
		if err != nil {
			return nil, err
		}
	default:
		break
	}

	return resp, nil
}

func (c *Service) throw(matrix string) {
	url := fmt.Sprintf("throw/%+v", matrix)
	_, err := c.sendToScoreboard(url, "post")
	if err != nil {
		log.Printf("Error when sending nextPlayer: %+v", err)
	}
}

func (c *Service) nextPlayer() {
	url := ("nextPlayer")
	_, err := c.sendToScoreboard(url, "post")
	if err != nil {
		log.Printf("Error when sending nextPlayer: %+v", err)
	}

	// Write 4 to serial to set bUltrasonicThresholdMeasured false
	c.Write("4")
}

func (c *Service) rematch() {
	url := ("rematch")
	_, err := c.sendToScoreboard(url, "post")
	if err != nil {
		log.Printf("Error when sending nextPlayer: %+v", err)
	}
}

func (c *Service) buttonOn() {
	c.Write("1")
}

func (c *Service) buttonOff() {
	c.Write("2")
}

func (c *Service) updateStatus() {
	log.Println("Updating game state")
	url := "display"
	resp, err := c.sendToScoreboard(url, "get")
	if err != nil {
		log.Println(err)
		return
	}
	json.NewDecoder(resp.Body).Decode(&c.State)
	log.Printf("game state in update function: %+v", c.State)
}
