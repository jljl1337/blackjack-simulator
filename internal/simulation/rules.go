package simulation

import (
	"github.com/jljl1337/blackjack-simulator/internal/blackjack"
)

type Rules interface {
	// GetActionsAllowed returns the actions available to the player in the current state of the game.
	GetActionsAllowed() (map[blackjack.Action]bool, error)
}

// PlayRules implements the rules of the game for a single round of blackjack.
type PlayRules struct{}

// NewPlayRules creates a new instance of PlayRules.
func NewPlayRules() PlayRules {
	return PlayRules{}
}

func (r PlayRules) GetActionsAllowed() (map[blackjack.Action]bool, error) {
	// This method should return the actions available to the player.
	return map[blackjack.Action]bool{
		blackjack.Hit:       true,
		blackjack.Stand:     true,
		blackjack.Double:    true,
		blackjack.Split:     true,
		blackjack.Surrender: false,
	}, nil
}
