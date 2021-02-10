package sender

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/dascr/dascr-machine/service/common"
	"github.com/dascr/dascr-machine/service/config"
	"github.com/dascr/dascr-machine/service/logger"
	"github.com/dascr/dascr-machine/service/serial"
	"github.com/dascr/dascr-machine/service/state"
)

// Sender will represent the connection between the scoreboard and the connector service
type Sender struct {
	User       string
	Pass       string
	HTTPS      bool
	IP         string
	Port       string
	GameID     string
	HTTPClient *http.Client
	Serial     *serial.Serial
}

// New will return an instantiated Sender
func New(serial *serial.Serial) *Sender {
	scoreboard := config.Config.Scoreboard

	sender := &Sender{
		User:   scoreboard.User,
		Pass:   scoreboard.Pass,
		HTTPS:  scoreboard.HTTPS,
		IP:     scoreboard.Host,
		Port:   scoreboard.Port,
		GameID: scoreboard.Game,
		HTTPClient: &http.Client{
			Jar: &cookiejar.Jar{},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				req.Header.Add("Authorization", "Basic "+common.BasicAuth(scoreboard.User, scoreboard.Pass))
				return nil
			},
		},
		Serial: serial,
	}
	return sender

}

// sendToScoreboard will send the httpo request to the scoreboard
// and update the state with the response received
func (s *Sender) sendToScoreboard(url, method string) error {
	var req *http.Request
	var resp *http.Response
	var err error

	protocol := "http"
	if s.HTTPS {
		protocol = "https"
	}

	target := fmt.Sprintf("%+v://%+v/api/game/%+v/%+v", protocol, s.IP, s.GameID, url)

	if s.Port != "80" && s.Port != "443" {
		target = fmt.Sprintf("%+v://%+v:%+v/api/game/%+v/%+v", protocol, s.IP, s.Port, s.GameID, url)
	}

	switch method {
	case "get":
		req, err = http.NewRequest("GET", target, nil)
		if err != nil {
			return err
		}
	case "post":
		req, err = http.NewRequest("POST", target, nil)
		if err != nil {
			return err
		}
	default:
		break
	}

	if s.User != "" {
		req.Header.Add("Authorization", "Basic "+common.BasicAuth(s.User, s.Pass))
	}

	resp, err = s.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	json.NewDecoder(resp.Body).Decode(&state.GameState)

	// Write state to Arduino
	switch state.GameState.GameState {
	case "THROW":
		s.Serial.Write("s,1")
	case "NEXTPLAYER":
		s.Serial.Write("s,2")
	case "BUST":
		s.Serial.Write("s,2")
	case "BUSTCONDITION":
		s.Serial.Write("s,2")
	case "BUSTNOCHECKOUT":
		s.Serial.Write("s,2")
	case "WON":
		s.Serial.Write("s,5")
	}

	return nil
}

// Throw will send a throw using sendToScoreboard
func (s *Sender) Throw(matrix string) {
	url := fmt.Sprintf("throw/%+v", matrix)
	err := s.sendToScoreboard(url, "post")
	if err != nil {
		logger.Errorf("Error when sending nextPlayer: %+v", err)
	}
}

// NextPlayer will send nextPlayer using sendToScoreboard
func (s *Sender) NextPlayer() {
	url := ("nextPlayer")
	err := s.sendToScoreboard(url, "post")
	if err != nil {
		logger.Errorf("Error when sending nextPlayer: %+v", err)
	}
}

// Rematch will send rematch using sendToScoreboard
func (s *Sender) Rematch() {
	url := ("rematch")
	err := s.sendToScoreboard(url, "post")
	if err != nil {
		logger.Errorf("Error when sending nextPlayer: %+v", err)
	}
}

// UpdateStatus will fetch and update the status from the scoreboard with sendToScoreboard
func (s *Sender) UpdateStatus() {
	url := "display"
	err := s.sendToScoreboard(url, "get")
	if err != nil {
		logger.Error(err)
	}
}

// CheckConnection will check if the scoreboard is reachable and return an error if not
// It is used by the ui to check after updating the Scoreboard config
func (s *Sender) CheckConnection() error {
	protocol := "http"
	if s.HTTPS {
		protocol = "https"
	}

	target := fmt.Sprintf("%+v://%+v/api", protocol, s.IP)
	if s.Port != "80" && s.Port != "443" {
		target = fmt.Sprintf("%+v://%+v:%+v/api", protocol, s.IP, s.Port)
	}

	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return err
	}

	if s.User != "" {
		req.Header.Add("Authorization", "Basic "+common.BasicAuth(s.User, s.Pass))
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		err := fmt.Errorf("Status is %+v, Connection not established", resp.StatusCode)
		return err
	}
	return nil
}
