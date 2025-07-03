package simulation

import (
	"github.com/jljl1337/blackjack-simulator/internal/blackjack"
)

// Rules implements the rules of the game for a single round of blackjack.
type Rules struct{}

// NewRules creates a new instance of PlayRules.
func NewRules() Rules {
	return Rules{}
}

func (r Rules) GetActionsAllowed(currentHandSize int, numHands int, splitAce bool) (map[blackjack.Action]bool, error) {
	// This method should return the actions available to the player.
	return map[blackjack.Action]bool{
		blackjack.Hit:       !splitAce,
		blackjack.Stand:     true,
		blackjack.Double:    !splitAce && currentHandSize == 2,
		blackjack.Split:     !splitAce && numHands < 4,
		blackjack.Surrender: true,
	}, nil
}
