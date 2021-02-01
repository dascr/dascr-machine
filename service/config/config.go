package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// default config
var (
	def = Settings{
		Machine: MachineConfig{
			WaitingTime: 5,
			Piezo:       20,
			Serial:      "/dev/ttyACM0",
		},
		Scoreboard: ScoreboardConfig{
			HTTPS: false,
			Host:  "127.0.0.1",
			Port:  "8000",
			Game:  "dascr",
			User:  "",
			Pass:  "",
		},
	}
	// Read config file
	name = "config.json"
)

// Config will hold the global configuration
var Config Settings

// Settings will hold the configuration as a main object
type Settings struct {
	Machine    MachineConfig    `json:"machine"`
	Scoreboard ScoreboardConfig `json:"scoreboard"`
}

// MachineConfig will hold the config for machine parameters
type MachineConfig struct {
	WaitingTime int    `json:"wait"`
	Piezo       int    `json:"piezo"`
	Serial      string `json:"serial"`
	Error       string `json:"-"`
}

// ScoreboardConfig will hold the config of the scoreboard
// to send to
type ScoreboardConfig struct {
	HTTPS bool   `json:"https"`
	Host  string `json:"host"`
	Port  string `json:"port"`
	Game  string `json:"game"`
	User  string `json:"user"`
	Pass  string `json:"pass"`
	Error string `json:"-"`
}

// Init will read config from file or create default config
func Init() error {
	_, err := os.Stat(name)
	if err != nil {
		// If file is not there write default config to file
		// and return using default config
		f, err := os.Create(name)
		if err != nil {
			return err
		}
		defer f.Close()

		Config = def

		err = SaveConfig()
		if err != nil {
			return err
		}

		return nil
	}

	// Otherwise read config from file and set
	// Read json from file
	jsonFile, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(jsonFile), &Config)
	if err != nil {
		return err
	}

	return nil
}

// SaveConfig will write the config to a file after updating it
func SaveConfig() error {
	// Write json to file
	output, err := json.MarshalIndent(Config, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(name, output, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
