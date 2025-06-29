package person

import "github.com/jljl1337/blackjack-simulator/internal/core"

type PlayerHand struct {
	Hand
	Bet int
}

func NewPlayerHand() *PlayerHand {
	return &PlayerHand{
		Hand: Hand{
			cards: []core.Card{},
		},
		Bet: 0,
	}
}
