package state

// State will represent the games state
type State struct {
	GameState string `json:"GameState"`
}

// GameState is the global var of the state
// used by multiple packages
var GameState State
