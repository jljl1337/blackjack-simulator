package simulation

import (
	"github.com/jljl1337/blackjack-simulator/internal/blackjack"
)

// Rules implements the rules of the game for a single round of blackjack.
type Rules struct {
	doubleAfterSplit    bool
	hitAfterSplitAce    bool
	splitAfterSplitAce  bool
	doubleAfterSplitAce bool
	maxNumHands         int
	surrenderAllowed    bool
}

// NewRules creates a new instance of PlayRules.
func NewRules(doubleAfterSplit, hitAfterSplitAce, splitAfterSplitAce, doubleAfterSplitAce bool, maxNumHands int, surrenderAllowed bool) Rules {
	return Rules{
		doubleAfterSplit:    doubleAfterSplit,
		hitAfterSplitAce:    hitAfterSplitAce,
		splitAfterSplitAce:  splitAfterSplitAce,
		doubleAfterSplitAce: doubleAfterSplitAce,
		maxNumHands:         maxNumHands,
		surrenderAllowed:    surrenderAllowed,
	}
}

func (r Rules) GetActionsAllowed(currentHandSize int, numHands int, splitAce bool) (map[blackjack.Action]bool, error) {
	// This method should return the actions available to the player.

	// Check if hit is allowed for split aces
	canHit := !splitAce || r.hitAfterSplitAce

	// Check if double is allowed
	canDouble := currentHandSize == 2 && (!splitAce || r.doubleAfterSplitAce) && (numHands < 2 || r.doubleAfterSplit)

	// Check if split is allowed
	canSplit := ((!splitAce || r.splitAfterSplitAce) && numHands < r.maxNumHands) || r.maxNumHands < 0

	return map[blackjack.Action]bool{
		blackjack.Hit:       canHit,
		blackjack.Stand:     true,
		blackjack.Double:    canDouble,
		blackjack.Split:     canSplit,
		blackjack.Surrender: r.surrenderAllowed,
	}, nil
}
