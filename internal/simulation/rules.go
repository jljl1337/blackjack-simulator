package simulation

import (
	"github.com/jljl1337/blackjack-simulator/internal/blackjack"
)

// Rules implements the rules of the game for a single round of blackjack.
type Rules struct{}

// NewPlayRules creates a new instance of PlayRules.
func NewPlayRules() Rules {
	return Rules{}
}

func (r Rules) GetActionsAllowed(currentHandSize int) (map[blackjack.Action]bool, error) {
	// This method should return the actions available to the player.
	return map[blackjack.Action]bool{
		blackjack.Hit:       true,
		blackjack.Stand:     true,
		blackjack.Double:    currentHandSize == 2,
		blackjack.Split:     true,
		blackjack.Surrender: true,
	}, nil
}
