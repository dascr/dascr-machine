package connector

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// basicAuth will add basic auth to connection
func (c *Service) basicAuth() string {
	auth := c.User + ":" + c.Pass
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *Service) redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+c.basicAuth())
	return nil
}

// sendToScoreboard will send the httpo request to the scoreboard
// and update the state with the response received
func (c *Service) sendToScoreboard(url, method string) error {
	var req *http.Request
	var resp *http.Response
	var err error

	protocol := "http"
	if c.HTTPS {
		protocol = "https"
	}

	target := fmt.Sprintf("%+v://%+v/api/game/%+v/%+v", protocol, c.Host, c.Game, url)

	if c.Port != "80" && c.Port != "443" {
		target = fmt.Sprintf("%+v://%+v:%+v/api/game/%+v/%+v", protocol, c.Host, c.Port, c.Game, url)
	}

	switch method {
	case "get":
		req, err = http.NewRequest("GET", target, nil)
		if err != nil {
			c.buttonBlink()
			return err
		}
	case "post":
		req, err = http.NewRequest("POST", target, nil)
		if err != nil {
			c.buttonBlink()
			return err
		}
	default:
		break
	}

	if c.User != "" {
		req.Header.Add("Authorization", "Basic "+c.basicAuth())
	}

	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		c.buttonBlink()
		return err
	}

	json.NewDecoder(resp.Body).Decode(&c.State)
	return nil
}

// throw will send a throw using sendToScoreboard
func (c *Service) throw(matrix string) {
	url := fmt.Sprintf("throw/%+v", matrix)
	err := c.sendToScoreboard(url, "post")
	if err != nil {
		log.Printf("Error when sending nextPlayer: %+v", err)
	}
}

// nextPlayer will send nextPlayer using sendToScoreboard
func (c *Service) nextPlayer() {
	url := ("nextPlayer")
	err := c.sendToScoreboard(url, "post")
	if err != nil {
		log.Printf("Error when sending nextPlayer: %+v", err)
	}

	// Write 4 to serial to set bUltrasonicThresholdMeasured false
	c.Write("4")
	c.State.Blink = false
}

// rematch will send rematch using sendToScoreboard
func (c *Service) rematch() {
	url := ("rematch")
	err := c.sendToScoreboard(url, "post")
	if err != nil {
		log.Printf("Error when sending nextPlayer: %+v", err)
	}
}

// updateStatus will fetch and update the status from the scoreboard with sendToScoreboard
func (c *Service) updateStatus() {
	url := "display"
	err := c.sendToScoreboard(url, "get")
	if err != nil {
		log.Println(err)
	}
}

// buttonOn will write 1 to serial and thus switch the button on
func (c *Service) buttonOn() {
	c.Write("1")
}

// buttonOff will write 2 to serial and thus switch the button off
func (c *Service) buttonOff() {
	c.Write("2")
}

// buttonBlink will write 7 to serial and thus let the button blink 7 times
func (c *Service) buttonBlink() {
	c.Write("7")
}
